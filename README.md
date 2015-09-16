This is an (working but unfinished) Golang implementation of:

- a reader for Blizzard's MPQ file format;
- a reader for Blizzard's Starcraft/Heroes of the Storm game replay formats;
- and various support tools for the above, including e.g. a code generator
  that translates the replay data layouts published by Blizzard into Golang
  decoders.

I had written it with the thought of doing some statistical analysis of
my replays, but in practice the replay format is more or less a recording
of user clicks, not of game logic, so it's difficult to extract anything
interesting (like, say, damage dealt) without a full implementation of all
the game logic (which itself varies per patch).
