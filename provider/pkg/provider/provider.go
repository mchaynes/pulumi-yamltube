// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/mchaynes/pulumi-yamltube/provider/pkg/internal/youtube"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/mapper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

type Playlist struct {
	PlaylistID  string   `json:"playlistId,omitempty"`
	ChannelID   string   `json:"channelId,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Visibility  string   `json:"visibility"`
	Videos      []string `json:"videos"`
}

type youtubeProvider struct {
	host    *provider.HostClient
	name    string
	schema  string
	version string
	youtube *youtube.YouTube
}

const (
	playlistType = "yamltube:youtube:Playlist"
)

func makeProvider(host *provider.HostClient, name, version, schema string) (pulumirpc.ResourceProviderServer, error) {
	// inject version into schema
	versionedSchema := mustSetSchemaVersion(schema, version)

	// Return the new provider
	return &youtubeProvider{
		host:    host,
		name:    name,
		schema:  versionedSchema,
		version: version,
	}, nil
}

func (k *youtubeProvider) Attach(context context.Context, req *pulumirpc.PluginAttach) (*empty.Empty, error) {
	host, err := provider.NewHostClient(req.GetAddress())
	if err != nil {
		return nil, err
	}
	k.host = host
	return &empty.Empty{}, nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (k *youtubeProvider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Call is not yet implemented")
}

// Construct creates a new component resource.
func (k *youtubeProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

// CheckConfig validates the configuration for this provider.
func (k *youtubeProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig diffs the configuration for this provider.
func (k *youtubeProvider) DiffConfig(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return &pulumirpc.DiffResponse{}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *youtubeProvider) Configure(ctx context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	y, err := youtube.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to youtube: %w", err)
	}
	k.youtube = y
	return &pulumirpc.ConfigureResponse{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *youtubeProvider) Invoke(ctx context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	tok := req.GetTok()
	return nil, fmt.Errorf("unknown Invoke token '%s'", tok)
}

// StreamInvoke dynamically executes a built-in function in the provider. The result is streamed
// back as a series of messages.
func (k *youtubeProvider) StreamInvoke(req *pulumirpc.InvokeRequest, server pulumirpc.ResourceProvider_StreamInvokeServer) error {
	tok := req.GetTok()
	return fmt.Errorf("unknown StreamInvoke token '%s'", tok)
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *youtubeProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != playlistType {
		return nil, fmt.Errorf("unknown resource type '%s'", ty)
	}
	news, err := UnmarshalPlaylist(req.GetNews())
	if err != nil {
		return nil, fmt.Errorf("check failed to unmarshal news: %w", err)
	}
	olds, err := UnmarshalPlaylist(req.GetOlds())
	if err != nil {
		return nil, err
	}
	news.PlaylistID = olds.PlaylistID
	news.ChannelID = olds.ChannelID
	newsMarshalled, err := MarshalPlaylist(news)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %w", err)
	}
	return &pulumirpc.CheckResponse{Inputs: newsMarshalled, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *youtubeProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != playlistType {
		return nil, fmt.Errorf("unknown resource type '%s'", ty)
	}
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		KeepUnknowns:     true,
		SkipNulls:        true,
		KeepOutputValues: true,
	})
	if err != nil {
		return nil, err
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		KeepUnknowns:     true,
		SkipNulls:        true,
		KeepOutputValues: true,
	})
	if err != nil {
		return nil, err
	}

	changes := pulumirpc.DiffResponse_DIFF_NONE
	var diffs, replaces []string
	properties := map[string]bool{
		"title":       false,
		"description": false,
		"visibility":  false,
		"videos":      false,
	}
	if d := olds.Diff(news); d != nil {
		for key, replace := range properties {
			i := sort.SearchStrings(req.IgnoreChanges, key)
			if i < len(req.IgnoreChanges) && req.IgnoreChanges[i] == key {
				continue
			}

			if d.Changed(resource.PropertyKey(key)) {
				changes = pulumirpc.DiffResponse_DIFF_SOME
				diffs = append(diffs, key)

				if replace {
					replaces = append(replaces, key)
				}
			}
		}
	}
	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Diffs:    diffs,
		Replaces: replaces,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (k *youtubeProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != playlistType {
		return nil, fmt.Errorf("unknown resource type '%s'", ty)
	}

	playlist, err := UnmarshalPlaylist(req.GetProperties())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	created, err := k.youtube.CreatePlaylist(ctx, playlist.Title, playlist.Description, playlist.Visibility)
	if err != nil {
		return nil, fmt.Errorf("failed to create playlist %q: %w", playlist.Title, err)
	}

	ids, err := youtube.ToVideoIds(playlist.Videos)
	if err != nil {
		return nil, err
	}
	_, err = k.youtube.SyncPlaylist(ctx, created.Id, ids)
	if err != nil {
		return nil, err
	}

	newPlaylist := Playlist{
		PlaylistID:  created.Id,
		Title:       created.Snippet.Title,
		Description: created.Snippet.Description,
		ChannelID:   created.Snippet.ChannelId,
		Visibility:  created.Status.PrivacyStatus,
		// this stays the same, since we do the mapping to ids on every run
		// and don't want to convert these into links
		Videos: playlist.Videos,
	}
	outputProperties, err := MarshalPlaylist(newPlaylist)
	if err != nil {
		return nil, err
	}
	return &pulumirpc.CreateResponse{
		Id:         newPlaylist.PlaylistID,
		Properties: outputProperties,
	}, nil
}

// Read the current live state associated with a resource.
func (k *youtubeProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != playlistType {
		return nil, fmt.Errorf("unknown resource type '%s'", ty)
	}
	return &pulumirpc.ReadResponse{
		Id:         req.GetId(),
		Inputs:     req.GetInputs(),
		Properties: req.GetProperties(),
	}, nil
}

// Update updates an existing resource with new values.
func (k *youtubeProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != playlistType {
		return nil, fmt.Errorf("unknown resource type '%s'", ty)
	}

	old, err := UnmarshalPlaylist(req.GetOlds())
	if err != nil {
		return nil, err
	}
	new, err := UnmarshalPlaylist(req.GetNews())
	if err != nil {
		return nil, err
	}
	new.PlaylistID = old.PlaylistID
	new.ChannelID = old.ChannelID
	if !reflect.DeepEqual(old, new) {
		_, err := k.youtube.UpdatePlaylist(ctx, new.PlaylistID, new.Title, new.Description, new.Visibility)
		if err != nil {
			return nil, fmt.Errorf("failed to update playlist: %w", err)
		}

		ids, err := youtube.ToVideoIds(new.Videos)
		if err != nil {
			return nil, fmt.Errorf("failed to parse video ids: %w", err)
		}
		_, err = k.youtube.SyncPlaylist(ctx, new.PlaylistID, ids)
		if err != nil {
			return nil, fmt.Errorf("failed to sync playlist videos: %w", err)
		}
	}

	outputProperties, err := MarshalPlaylist(new)
	if err != nil {
		return nil, err
	}
	return &pulumirpc.UpdateResponse{
		Properties: outputProperties,
	}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *youtubeProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*empty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != playlistType {
		return nil, fmt.Errorf("unknown resource type '%s'", ty)
	}
	var playlist Playlist
	playlist, err := UnmarshalPlaylist(req.GetProperties())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize playlist properties: %w", err)
	}
	err = k.youtube.DeletePlaylist(ctx, playlist.PlaylistID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete playlist: %w", err)
	}

	return &empty.Empty{}, nil
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *youtubeProvider) GetPluginInfo(context.Context, *empty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}

// GetSchema returns the JSON-serialized schema for the provider.
func (k *youtubeProvider) GetSchema(ctx context.Context, req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	return &pulumirpc.GetSchemaResponse{
		Schema: k.schema,
	}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (k *youtubeProvider) Cancel(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// mustSetSchemaVersion deserializes schemaStr from json, sets Version field
// then serializes back to json string
func mustSetSchemaVersion(schemaStr string, version string) string {
	var spec schema.PackageSpec
	if err := json.Unmarshal([]byte(schemaStr), &spec); err != nil {
		panic(fmt.Errorf("failed to parse schema: %w", err))
	}
	spec.Version = version
	bytes, err := json.Marshal(spec)
	if err != nil {
		panic(fmt.Errorf("failed to serialize versioned schema: %w", err))
	}
	return string(bytes)
}

func UnmarshalPlaylist(props *structpb.Struct) (Playlist, error) {
	inputProps, err := plugin.UnmarshalProperties(props, plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: false})
	if err != nil {
		return Playlist{}, err
	}
	inputs := inputProps.Mappable()
	var out Playlist
	err = mapper.MapIM(inputs, &out)
	if err != nil {
		return out, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return out, nil
}

func MarshalPlaylist(p Playlist) (*structpb.Struct, error) {
	outputs, err := mapper.New(&mapper.Opts{
		IgnoreMissing:      true,
		IgnoreUnrecognized: true,
	}).Encode(p)
	if err != nil {
		return nil, err
	}

	return plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputs),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true},
	)
}
