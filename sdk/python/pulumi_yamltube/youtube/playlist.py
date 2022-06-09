# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from .. import _utilities

__all__ = ['PlaylistArgs', 'Playlist']

@pulumi.input_type
class PlaylistArgs:
    def __init__(__self__, *,
                 description: pulumi.Input[str],
                 title: pulumi.Input[str],
                 videos: pulumi.Input[Sequence[pulumi.Input[str]]],
                 visibility: pulumi.Input[str],
                 channel_id: Optional[pulumi.Input[str]] = None,
                 playlist_id: Optional[pulumi.Input[str]] = None):
        """
        The set of arguments for constructing a Playlist resource.
        :param pulumi.Input[str] description: description of the playlist
        :param pulumi.Input[str] title: title of the playlist
        :param pulumi.Input[Sequence[pulumi.Input[str]]] videos: array of youtube links or ids. (https://www.youtube.com/watch?v=dQw4w9WgXcQ, or dQw4w9WgXcQ)
        :param pulumi.Input[str] visibility: visibility of the playlist. valid values are public, private, unlisted
        :param pulumi.Input[str] channel_id: the id of the channel. https://www.youtube.com/channel/<channelId>
        """
        pulumi.set(__self__, "description", description)
        pulumi.set(__self__, "title", title)
        pulumi.set(__self__, "videos", videos)
        pulumi.set(__self__, "visibility", visibility)
        if channel_id is not None:
            pulumi.set(__self__, "channel_id", channel_id)
        if playlist_id is not None:
            pulumi.set(__self__, "playlist_id", playlist_id)

    @property
    @pulumi.getter
    def description(self) -> pulumi.Input[str]:
        """
        description of the playlist
        """
        return pulumi.get(self, "description")

    @description.setter
    def description(self, value: pulumi.Input[str]):
        pulumi.set(self, "description", value)

    @property
    @pulumi.getter
    def title(self) -> pulumi.Input[str]:
        """
        title of the playlist
        """
        return pulumi.get(self, "title")

    @title.setter
    def title(self, value: pulumi.Input[str]):
        pulumi.set(self, "title", value)

    @property
    @pulumi.getter
    def videos(self) -> pulumi.Input[Sequence[pulumi.Input[str]]]:
        """
        array of youtube links or ids. (https://www.youtube.com/watch?v=dQw4w9WgXcQ, or dQw4w9WgXcQ)
        """
        return pulumi.get(self, "videos")

    @videos.setter
    def videos(self, value: pulumi.Input[Sequence[pulumi.Input[str]]]):
        pulumi.set(self, "videos", value)

    @property
    @pulumi.getter
    def visibility(self) -> pulumi.Input[str]:
        """
        visibility of the playlist. valid values are public, private, unlisted
        """
        return pulumi.get(self, "visibility")

    @visibility.setter
    def visibility(self, value: pulumi.Input[str]):
        pulumi.set(self, "visibility", value)

    @property
    @pulumi.getter(name="channelId")
    def channel_id(self) -> Optional[pulumi.Input[str]]:
        """
        the id of the channel. https://www.youtube.com/channel/<channelId>
        """
        return pulumi.get(self, "channel_id")

    @channel_id.setter
    def channel_id(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "channel_id", value)

    @property
    @pulumi.getter(name="playlistId")
    def playlist_id(self) -> Optional[pulumi.Input[str]]:
        return pulumi.get(self, "playlist_id")

    @playlist_id.setter
    def playlist_id(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "playlist_id", value)


class Playlist(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 channel_id: Optional[pulumi.Input[str]] = None,
                 description: Optional[pulumi.Input[str]] = None,
                 playlist_id: Optional[pulumi.Input[str]] = None,
                 title: Optional[pulumi.Input[str]] = None,
                 videos: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 visibility: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        """
        Create a Playlist resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[str] channel_id: the id of the channel. https://www.youtube.com/channel/<channelId>
        :param pulumi.Input[str] description: description of the playlist
        :param pulumi.Input[str] title: title of the playlist
        :param pulumi.Input[Sequence[pulumi.Input[str]]] videos: array of youtube links or ids. (https://www.youtube.com/watch?v=dQw4w9WgXcQ, or dQw4w9WgXcQ)
        :param pulumi.Input[str] visibility: visibility of the playlist. valid values are public, private, unlisted
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: PlaylistArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Create a Playlist resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param PlaylistArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(PlaylistArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 channel_id: Optional[pulumi.Input[str]] = None,
                 description: Optional[pulumi.Input[str]] = None,
                 playlist_id: Optional[pulumi.Input[str]] = None,
                 title: Optional[pulumi.Input[str]] = None,
                 videos: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 visibility: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        if opts is None:
            opts = pulumi.ResourceOptions()
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.version is None:
            opts.version = _utilities.get_version()
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = PlaylistArgs.__new__(PlaylistArgs)

            __props__.__dict__["channel_id"] = channel_id
            if description is None and not opts.urn:
                raise TypeError("Missing required property 'description'")
            __props__.__dict__["description"] = description
            __props__.__dict__["playlist_id"] = playlist_id
            if title is None and not opts.urn:
                raise TypeError("Missing required property 'title'")
            __props__.__dict__["title"] = title
            if videos is None and not opts.urn:
                raise TypeError("Missing required property 'videos'")
            __props__.__dict__["videos"] = videos
            if visibility is None and not opts.urn:
                raise TypeError("Missing required property 'visibility'")
            __props__.__dict__["visibility"] = visibility
        super(Playlist, __self__).__init__(
            'yamltube:youtube:Playlist',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'Playlist':
        """
        Get an existing Playlist resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = PlaylistArgs.__new__(PlaylistArgs)

        __props__.__dict__["channel_id"] = None
        __props__.__dict__["description"] = None
        __props__.__dict__["playlist_id"] = None
        __props__.__dict__["title"] = None
        __props__.__dict__["videos"] = None
        __props__.__dict__["visibility"] = None
        return Playlist(resource_name, opts=opts, __props__=__props__)

    @property
    @pulumi.getter(name="channelId")
    def channel_id(self) -> pulumi.Output[str]:
        """
        the id of the channel. https://www.youtube.com/channel/<channelId>
        """
        return pulumi.get(self, "channel_id")

    @property
    @pulumi.getter
    def description(self) -> pulumi.Output[str]:
        """
        description of the playlist
        """
        return pulumi.get(self, "description")

    @property
    @pulumi.getter(name="playlistId")
    def playlist_id(self) -> pulumi.Output[str]:
        return pulumi.get(self, "playlist_id")

    @property
    @pulumi.getter
    def title(self) -> pulumi.Output[str]:
        """
        title of the playlist
        """
        return pulumi.get(self, "title")

    @property
    @pulumi.getter
    def videos(self) -> pulumi.Output[Sequence[str]]:
        """
        array of youtube links or ids. (https://www.youtube.com/watch?v=dQw4w9WgXcQ, or dQw4w9WgXcQ)
        """
        return pulumi.get(self, "videos")

    @property
    @pulumi.getter
    def visibility(self) -> pulumi.Output[str]:
        """
        visibility of the playlist. valid values are public, private, unlisted
        """
        return pulumi.get(self, "visibility")
