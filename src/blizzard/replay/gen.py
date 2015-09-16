from s2protocol import protocol34835 as proto

# The typeinfo id for an empty struct, handled specially.
EMPTY_STRUCT_ID = 78


def fieldToGo(name):
    if name.startswith('m_'):
        name = name[2:]
    return name[0].upper() + name[1:]


def genReadInt(offset, bits):
    if offset == 0:
        return 'readBits(r, %(bits)d)' % locals()
    elif bits == 0:
        return '%d' % offset
    else:
        return '%(offset)d + int64(readBits(r, %(bits)d))' % locals()

def simple_name(name):
    for prefix in ['NNet.Game.S', 'NNet.Replay.Tracker.S']:
        if name.startswith(prefix):
            return name[len(prefix):]
    return name

print '''package replay

import (
"fmt"
"io"
)
'''

class TypeInfo(object):
    def __init__(self, id):
        self.id = id
        self.name = None
        self.names = set()
        self.name_hints = set()
        self.typ = None
        self.event = False
    
    def decodeCode(self, name=None):
        if self.kind == 'int':
            bounds, = self.args
            offset, bits = bounds
            return '%s(%s)' % (self.typ, genReadInt(offset, bits))
        elif self.kind == 'bool':
            return 'readBits(r, 1) != 0'
        elif self.kind == 'null':
            return 'nil'
        elif self.id == EMPTY_STRUCT_ID:
            return '&%s{}' % name
        return 'decode%s(r)' % (name or self.name)

typeinfos = [TypeInfo(i) for i in range(len(proto.typeinfos))]

# Name all the TypeInfos we have names for.
all_event_types = []
all_event_types += proto.game_event_types.values()
all_event_types += proto.message_event_types.values()
all_event_types += proto.tracker_event_types.values()
for i, name in all_event_types:
    name = simple_name(name)
    typeinfos[i].names.add(name)
    typeinfos[i].event = True

typeinfos[7].name = 'GameLoopDeltaAuto'
typeinfos[7].typ = 'int32'
typeinfos[18].name = 'Header'
typeinfos[80].name = 'CameraTarget'

# Translate proto.typeinfos into TypeInfos.
for i, (f, args) in enumerate(proto.typeinfos):
    typeinfo = typeinfos[i]
    assert f.startswith('_')
    f = f[1:]
    typeinfo.kind = f
    typeinfo.args = args
    if f == 'int':
        bounds, = args
        offset, bits = bounds
        if not typeinfo.name:
            if offset == 0:
                if bits <= 8:
                    typeinfo.typ = 'int8'
                elif bits <= 16:
                    typeinfo.typ = 'int16'
                elif bits <= 32:
                    typeinfo.typ = 'int32'
                else:
                    typeinfo.typ = 'int64'
            else:
                typeinfo.name = 'Int%d' % bits
                if offset != 0:
                    typeinfo.name += '_%d' % i
        if typeinfo.typ is None:
            typeinfo.typ = 'int64'
    elif f == 'choice':
        if not typeinfo.typ:
            typeinfo.typ = 'interface{}'
    elif f == 'optional':
        ref, = args
        typeinfo.ref = ref
    elif f == 'blob':
        bounds, = args
        offset, bits = bounds
        typeinfo.name = 'ByteString_%d_%d' % (offset, bits)
        typeinfo.typ = 'string'
    elif f == 'struct':
        fields, = args
        typeinfo.fields = fields
        for field, id, x in fields:
            # TODO: x is some sort of field tag id.
            name = fieldToGo(field)
            typeinfos[id].name_hints.add(name)
    elif f == 'array':
        bounds, ref = args
        typeinfo.ref = ref
    elif f == 'bool':
        typeinfo.typ = 'bool'
    elif f == 'bitarray':
        typeinfo.typ = 'uint64'
    elif f in ('fourcc', 'bitarray'):
        typeinfo.typ = 'TODO'

# Try to assign names and types to TypeInfos.
for ti in typeinfos:
    if ti.kind == 'struct':
        if not ti.name:
            if len(ti.names) > 0:
                ti.name = list(ti.names)[0]
            elif len(ti.name_hints) > 0:
                ti.name = list(ti.name_hints)[0]
            else:
                ti.name = 'Unknown%d' % ti.id
        ti.typ = '*' + ti.name
    elif ti.kind == 'optional':
        ti.typ = '*%s' % typeinfos[ti.ref].typ
    elif ti.kind == 'array':
        ti.typ = '[]%s' % typeinfos[ti.ref].typ

    if ti.name is None:
        ti.name = 'Unknown%d' % ti.id

print """type EventMeta struct {
GameLoop int
UserId int
}
type Event interface {
  Meta() *EventMeta
}
func (e *EventMeta) Meta() *EventMeta {
return e
}
"""

for ti in typeinfos:
    print
    if ti.typ == 'TODO':
        print '// TODO: %s (%d)' % (ti.name, ti.id)
        continue

    print '// typeinfo %d (%s)' % (ti.id, ti.kind)
    if ti.kind == 'struct':
        if len(ti.names) > 1:
            print '// names: %r' % ti.names
        if not ti.names and len(ti.name_hints) > 1:
            print '// name hints: %r' % ti.name_hints

        names = ti.names
        if not names and len(ti.name_hints) > 1:
            names = list(ti.name_hints)[:1]
        if not names:
            names = [ti.name]

        for name in names:
            print 'type %(name)s struct {' % locals()
            if ti.event:
                print 'EventMeta'
            for field, id, x in ti.fields:
                # TODO: x is some sort of field tag id.
                name = fieldToGo(field)
                typ = typeinfos[id].typ
                if typ == 'TODO':
                    print '// TODO:',
                print '%(name)s %(typ)s // %(id)d' % locals()
            print '}'

    if ti.kind in ('null', 'int', 'bool'):
        continue
    if ti.id == EMPTY_STRUCT_ID:
        continue

    for name in (ti.names or [ti.name]):
        typ = ti.typ
        if len(ti.names) > 0:
            typ = '*' + name
        print 'func decode%(name)s(r *BitReader) %(typ)s {' % locals()

        if ti.kind == 'int':
            bounds, = ti.args
            offset, bits = bounds
            print 'return %s(%s)' % (ti.typ, genReadInt(offset, bits))
        elif ti.kind == 'choice':
            tagbounds, values = ti.args
            print 'switch tag := %s; tag {' % genReadInt(*tagbounds)
            for tag in sorted(values.keys()):
                name, typ = values[tag]
                tti = typeinfos[typ]
                decode = tti.decodeCode()
                print 'case %(tag)d:  // %(name)s' % locals()
                if ti.typ == 'interface{}':
                    print 'return %s' % decode
                else:
                    print 'return %s(%s)' % (ti.typ, decode)
            print 'default: panic(fmt.Errorf("unknown choice tag %d", tag))'
            print '}'
        elif ti.kind == 'struct':
            fields, = ti.args
            print 'out := &%(name)s{}' % locals()
            for field, typ, x in fields:
                # TODO: x is some sort of field tag id.
                fieldname = fieldToGo(field)
                fti = typeinfos[typ]
                decode = fti.decodeCode()
                if fti.typ == 'TODO':
                    print 'panic("TODO")  // decode %(fieldname)s' % locals()
                else:
                    print 'out.%(fieldname)s = %(decode)s' % locals()
            print 'return out'
        elif ti.kind == 'blob':
            bounds, = ti.args
            offset, bits = bounds
            readInt = genReadInt(offset, bits)
            print '''n := %(readInt)s
            r.SyncToByte()
            buf := make([]byte, n)
            _, err := io.ReadFull(r, buf)
            if err != nil { panic(err) }
            return string(buf)''' % locals()
        elif ti.kind == 'array':
            bounds, typ = ti.args
            offset, bits = bounds
            elemti = typeinfos[typ]
            elemtyp = elemti.typ
            decode = elemti.decodeCode()
            readInt = genReadInt(offset, bits)
            print '''n := int(%(readInt)s)
            arr := make([]%(elemtyp)s, n)
            for i := 0; i < n; i++ {
            arr[i] = %(decode)s
            }
            return arr''' % locals()
        elif ti.kind == 'optional':
            typ, = ti.args
            argti = typeinfos[typ]
            decode = argti.decodeCode()
            print '''if readBits(r, 1) != 0 {
            ret := %(decode)s
            return &ret
            }
            return nil''' % locals()
        elif ti.kind == 'fourcc':
            print 'panic("TODO")'
        elif ti.kind == 'bitarray':
            print 'panic("TODO")'
        elif ti.kind == 'null':
            print 'panic("TODO")'
        else:
            assert False, f
        print '}'

for name, event_types in [('Game', proto.game_event_types),
                          ('Message', proto.message_event_types),
                          ('Tracker', proto.tracker_event_types)]:
  print '''func read%(name)sEvent(r *BitReader, typ int) Event {
  switch typ {''' % locals()
  for typ, (i, name) in event_types.iteritems():
      name = simple_name(name)
      print 'case %d: return %s' % (typ, typeinfos[i].decodeCode(name))
  print '''default:
  panic(fmt.Errorf("unknown event type %d", typ))
  }
  }'''
