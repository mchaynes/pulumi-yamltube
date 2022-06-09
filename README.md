# YamlTube's Pulumi Provider

This repo contains a pulumi native provider for yamltube. It contacts and synchronizes YouTube playlists on your behalf.

If you want to use YamlTube, go fork the [YamlTube repo](https://github.com/mchaynes/yamltube)

### What does a YouTube playlist in Yaml look like?

```yaml
name: yaml-rickroll
runtime: yaml
description: A rick roll playlist
resources:
  rickroll:
    type: yamltube:youtube:Playlist
    properties:
      title: "Rick Roll"
      description: "I couldn't think of a better example"
      visibility: public
      videos:
        - https://www.youtube.com/watch?v=dQw4w9WgXcQ

outputs:
  # output a link to the playlist
  playlist: https://www.youtube.com/playlist?list=${rickroll.playlistId}
```

### Does this support other Pulumi supported languages?

Not yet, because I'm lazy about getting publishing set up. I'll do it soon™
