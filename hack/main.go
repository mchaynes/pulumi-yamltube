package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	client := getClient(youtube.YoutubeScope)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create youtube service: %w", err)
	}
	/* FETCH PLAYLIST */
	// if err := getPlaylistAndItems(service); err != nil {
	// 	return err
	// }
	/* FETCH PLAYLIST */

	/* SYNC PLAYLIST */
	// playlistId := "PLeQFt2AXw9mSQpqcBfHkufqpBsS2x4hTD"
	// videoIds := []string{"BdEe5SpdIuo", "aqlGlaNlcWE"}
	// if _, err := syncPlaylist(ctx, service, playlistId, videoIds); err != nil {
	// 	return fmt.Errorf("failed to sync playlist: %w", err)
	// }
	/* SYNC PLAYLIST */

	/* CREATE PLAYLIST */
	playlist, err := CreatePlaylist(ctx, service, "YamlTubeTest", "Some Description", "public")
	if err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}
	mustPrettyPrint(playlist.MarshalJSON)
	/* CREATE PLAYLIST */

	return nil
}

func getPlaylistAndItems(service *youtube.Service) error {
	playlistsCall := service.Playlists.List([]string{"snippet,contentDetails"})
	playlistsCall.Mine(true)
	playlistsCall.MaxResults(100)
	playlistsResp, err := playlistsCall.Do()
	if err != nil {
		return fmt.Errorf("faild to get playlists: %w", err)
	}
	mustPrettyPrint(playlistsResp.MarshalJSON)

	// itemsListCall := service.PlaylistItems.List([]string{"snippet,contentDetails"})
	// itemsListCall.PlaylistId("PLeQFt2AXw9mTdc30PLoUNfyqJPc3cCgnM")
	// items, err := itemsListCall.Do()
	// if err != nil {
	// 	return fmt.Errorf("failed to get playlist items: %w", err)
	// }
	// mustPrettyPrint(items.MarshalJSON)

	return nil
}

type PlaylistInsert struct {
	VideoId  string
	Position int64
}

type PlaylistItemDelete struct {
	ItemId string
}

type PlaylistDiffResult struct {
	Inserts []PlaylistInsert
	Deletes []PlaylistItemDelete
}

func diffPlaylist(wantIds []string, gotItems []*youtube.PlaylistItem) PlaylistDiffResult {
	diff := PlaylistDiffResult{}
	// if what we got is longer than what we want, remove everything that's after
	// what we want
	if len(gotItems) > len(wantIds) {
		for i := len(wantIds); i < len(gotItems); i++ {
			diff.Deletes = append(diff.Deletes, PlaylistItemDelete{
				ItemId: gotItems[i].Id,
			})
		}
	}
	for i, wantVideoId := range wantIds {
		// check if we're out of bounds
		if len(gotItems) > i {
			gotItem := gotItems[i]
			// happy case, nothing to do
			if wantVideoId == gotItem.ContentDetails.VideoId {
				continue
			}
			// delete the playlistItem since it doesnt match
			diff.Deletes = append(diff.Deletes, PlaylistItemDelete{
				ItemId: gotItem.Id,
			})

			// we fall through to the insert here because
			// we didn't find the playlistItem we were looking for,
			// so we stil need to add it to the playlist
		}
		diff.Inserts = append(diff.Inserts, PlaylistInsert{
			VideoId:  wantVideoId,
			Position: int64(i),
		})
	}
	return diff
}

func syncPlaylist(ctx context.Context, service *youtube.Service, playlistId string, wantIds []string) (*PlaylistDiffResult, error) {
	items, err := fetchPlaylistItems(ctx, service, playlistId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlist for sync: %w", err)
	}

	diff := diffPlaylist(wantIds, items)

	for _, item := range diff.Deletes {
		err := service.PlaylistItems.Delete(item.ItemId).
			Context(ctx).
			Do()
		if err != nil {
			return nil, fmt.Errorf("failed to delete playlistItem %q: %w", item, err)
		}
	}

	for _, video := range diff.Inserts {
		_, err := service.PlaylistItems.Insert([]string{"id,snippet"}, &youtube.PlaylistItem{
			Snippet: &youtube.PlaylistItemSnippet{
				PlaylistId: playlistId,
				Position:   video.Position,
				ResourceId: &youtube.ResourceId{
					Kind:    "youtube#video",
					VideoId: video.VideoId,
				},
			},
		}).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("failed to insert playlistItem: %w", err)
		}
	}
	return &diff, nil
}

func fetchPlaylistItems(ctx context.Context, service *youtube.Service, playlistId string) ([]*youtube.PlaylistItem, error) {
	var pageToken string
	var items []*youtube.PlaylistItem
	firstLoop := true
	for len(pageToken) > 0 || firstLoop {
		itemsListCall := service.PlaylistItems.List([]string{"snippet,contentDetails"})
		itemsListCall.PlaylistId(playlistId)
		itemsListCall.Context(ctx)
		itemsListCall.MaxResults(50)
		itemsListCall.PageToken(pageToken)

		// execute
		listResp, err := itemsListCall.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch items: %w", err)
		}
		pageToken = listResp.NextPageToken
		items = append(items, listResp.Items...)
		firstLoop = false
	}
	return items, nil
}

func CreatePlaylist(ctx context.Context, service *youtube.Service, title, desc, visibility string) (*youtube.Playlist, error) {
	insert := service.Playlists.Insert([]string{"snippet,status"}, &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:       title,
			Description: desc,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: visibility,
		},
	})
	playlist, err := insert.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}
	return playlist, nil
}

func mustPrettyPrint(f func() ([]byte, error)) {
	data, err := f()
	if err != nil {
		panic(fmt.Errorf("failed to marshal data to json: %w", err))
	}
	buf := bytes.NewBuffer([]byte{})
	json.Indent(buf, data, " ", " ")
	fmt.Println(buf.String())
}

func listToMap[T comparable](list []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, v := range list {
		m[v] = struct{}{}
	}
	return m
}
