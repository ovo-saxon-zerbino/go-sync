//go:build !integration

package usergroup

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gosync "github.com/ovotech/go-sync"
)

func TestUserGroup_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	slackClient := newMockISlackUserGroup(t)

	adapter := &UserGroup{
		client:      slackClient,
		userGroupID: "test",
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
	}

	slackClient.EXPECT().GetUserGroupMembersContext(ctx, "test").Return([]string{"foo", "bar"}, nil)
	slackClient.EXPECT().GetUsersInfoContext(ctx, "foo", "bar").Maybe().Return(&[]slack.User{
		{ID: "foo", Profile: slack.UserProfile{Email: "foo@email"}},
		{ID: "bar", Profile: slack.UserProfile{Email: "bar@email"}},
	}, nil)
	slackClient.EXPECT().GetUsersInfoContext(ctx, "bar", "foo").Maybe().Return(&[]slack.User{
		{ID: "bar", Profile: slack.UserProfile{Email: "bar@email"}},
		{ID: "foo", Profile: slack.UserProfile{Email: "foo@email"}},
	}, nil)

	users, err := adapter.Get(ctx)

	require.NoError(t, err)
	assert.Equal(t, []string{"foo@email", "bar@email"}, users)
}

func TestUserGroup_Get_Pagination(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	slackClient := newMockISlackUserGroup(t)

	adapter := &UserGroup{
		client:      slackClient,
		userGroupID: "test",
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
	}

	incrementingSlice := make([]string, 60)
	firstPage := make([]interface{}, 30)
	secondPage := make([]interface{}, 30)
	firstResponse := make([]slack.User, 30)
	secondResponse := make([]slack.User, 30)

	for idx := range incrementingSlice {
		incrementingSlice[idx] = strconv.Itoa(idx)

		if idx < 30 {
			firstPage[idx] = strconv.Itoa(idx)
			firstResponse[idx] = slack.User{
				ID: strconv.Itoa(idx), IsBot: false, Profile: slack.UserProfile{Email: strconv.Itoa(idx)},
			}
		} else {
			secondPage[idx-30] = strconv.Itoa(idx)
			secondResponse[idx-30] = slack.User{
				ID: strconv.Itoa(idx), IsBot: false, Profile: slack.UserProfile{Email: strconv.Itoa(idx)},
			}
		}
	}

	slackClient.EXPECT().GetUserGroupMembersContext(ctx, "test").Return(incrementingSlice, nil)

	slackClient.EXPECT().GetUsersInfoContext(ctx, firstPage...).Return(&firstResponse, nil)
	slackClient.EXPECT().GetUsersInfoContext(ctx, secondPage...).Return(&secondResponse, nil)

	_, err := adapter.Get(ctx)

	require.NoError(t, err)
}

func TestUserGroup_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("No cache", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)

		adapter := &UserGroup{
			client:      slackClient,
			userGroupID: "test",
			Logger:      log.New(os.Stdout, "", log.LstdFlags),
		}

		err := adapter.Add(ctx, []string{"foo", "bar"})

		require.Error(t, err)
		require.ErrorIs(t, err, gosync.ErrCacheEmpty)
	})

	t.Run("Add accounts", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)

		adapter := &UserGroup{
			client:      slackClient,
			userGroupID: "test",
			Logger:      log.New(os.Stdout, "", log.LstdFlags),
		}

		slackClient.EXPECT().GetUserByEmailContext(ctx, "fizz@email").Return(&slack.User{ID: "fizz"}, nil)
		slackClient.EXPECT().GetUserByEmailContext(ctx, "buzz@email").Return(&slack.User{ID: "buzz"}, nil)
		slackClient.EXPECT().UpdateUserGroupMembersContext(ctx,
			"test", mock.Anything,
		).Run(func(_ context.Context, userGroup string, members string) { //nolint:contextcheck
			assert.Equal(t, "test", userGroup)
			assert.ElementsMatch(t, strings.Split(members, ","), []string{"foo", "bar", "fizz", "buzz"})
		}).Return(slack.UserGroup{DateDelete: 0}, nil)

		adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}
		err := adapter.Add(ctx, []string{"fizz@email", "buzz@email"})

		require.NoError(t, err)
		assert.Contains(t, adapter.cache, "fizz@email")
		assert.Contains(t, adapter.cache, "buzz@email")
	})
}

func TestUserGroup_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("No cache", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)

		adapter := &UserGroup{
			client:      slackClient,
			userGroupID: "test",
			Logger:      log.New(os.Stdout, "", log.LstdFlags),
		}

		err := adapter.Remove(ctx, []string{"foo@email"})

		require.Error(t, err)
		require.ErrorIs(t, err, gosync.ErrCacheEmpty)
	})

	t.Run("Remove accounts", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)

		adapter := &UserGroup{
			client:      slackClient,
			userGroupID: "test",
			cache:       map[string]string{"foo@email": "foo", "bar@email": "bar"},
			Logger:      log.New(os.Stdout, "", log.LstdFlags),
		}

		slackClient.EXPECT().UpdateUserGroupMembersContext(ctx, "test", "foo").Return(slack.UserGroup{}, nil)

		err := adapter.Remove(ctx, []string{"bar@email"})

		require.NoError(t, err)
		assert.Contains(t, adapter.cache, "foo@email")
		assert.NotContains(t, adapter.cache, "bar@email")
	})

	t.Run("Return/mute error if number of accounts reaches zero", func(t *testing.T) {
		t.Parallel()

		// Mock the error returned from the Slack API.
		errInvalidArguments := errors.New("invalid_arguments") //nolint:goerr113

		slackClient := newMockISlackUserGroup(t)

		adapter := &UserGroup{
			client:                 slackClient,
			userGroupID:            "test",
			cache:                  map[string]string{"foo@email": "foo"},
			MuteGroupCannotBeEmpty: false,
			Logger:                 log.New(os.Stdout, "", log.LstdFlags),
		}

		slackClient.EXPECT().UpdateUserGroupMembersContext(ctx, "test", "").Return(slack.UserGroup{}, errInvalidArguments)

		err := adapter.Remove(ctx, []string{"foo@email"})

		require.ErrorIs(t, err, errInvalidArguments)

		// Reset the cache and mute the empty group error.
		adapter.MuteGroupCannotBeEmpty = true

		err = adapter.Remove(ctx, []string{"foo@email"})

		require.NoError(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			SlackAPIKey: "test",
			UserGroupID: "usergroup",
		})

		require.NoError(t, err)
		assert.IsType(t, &UserGroup{}, adapter)
		assert.Equal(t, "usergroup", adapter.userGroupID)
		assert.False(t, adapter.MuteGroupCannotBeEmpty)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				UserGroupID: "usergroup",
			})

			require.ErrorIs(t, err, gosync.ErrMissingConfig)
			require.ErrorContains(t, err, SlackAPIKey)
		})

		t.Run("missing name", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				SlackAPIKey: "test",
			})

			require.ErrorIs(t, err, gosync.ErrMissingConfig)
			require.ErrorContains(t, err, UserGroupID)
		})
	})

	t.Run("MuteRestrictedErrOnKickFromPublic", func(t *testing.T) {
		t.Parallel()

		for _, test := range []string{"", "false", "FALSE", "False", "foobar", "test"} {
			adapter, err := Init(ctx, map[gosync.ConfigKey]string{
				SlackAPIKey:            "test",
				UserGroupID:            "usergroup",
				MuteGroupCannotBeEmpty: test,
			})

			require.NoError(t, err)
			assert.False(t, adapter.MuteGroupCannotBeEmpty, test)
		}

		for _, test := range []string{"true", "True", "TRUE"} {
			adapter, err := Init(ctx, map[gosync.ConfigKey]string{
				SlackAPIKey:            "test",
				UserGroupID:            "usergroup",
				MuteGroupCannotBeEmpty: test,
			})

			require.NoError(t, err)
			assert.True(t, adapter.MuteGroupCannotBeEmpty, test)
		}
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(os.Stderr, "custom logger", log.LstdFlags)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			SlackAPIKey: "test",
			UserGroupID: "usergroup",
		}, WithLogger(logger))

		require.NoError(t, err)
		assert.Equal(t, logger, adapter.Logger)
	})

	t.Run("with client", func(t *testing.T) {
		t.Parallel()

		client := slack.New("test")

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			UserGroupID: "usergroup",
		}, WithClient(client))

		require.NoError(t, err)
		assert.Equal(t, client, adapter.client)
	})
}
