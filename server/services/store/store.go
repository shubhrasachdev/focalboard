//go:generate mockgen --build_flags=--mod=mod -destination=mockstore/mockstore.go -package mockstore . Store
//go:generate go run ./generators/main.go
package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/mattermost/focalboard/server/model"
)

// Store represents the abstraction of the data storage.
type Store interface {
	GetBlocksWithParentAndType(boardID, parentID string, blockType string) ([]model.Block, error)
	GetBlocksWithParent(boardID, parentID string) ([]model.Block, error)
	GetBlocksWithRootID(boardID, rootID string) ([]model.Block, error)
	GetBlocksWithType(boardID, blockType string) ([]model.Block, error)
	GetSubTree2(boardID, blockID string, opts model.QuerySubtreeOptions) ([]model.Block, error)
	GetSubTree3(boardID, blockID string, opts model.QuerySubtreeOptions) ([]model.Block, error)
	GetBlocksForBoard(boardID string) ([]model.Block, error)
	// @withTransaction
	InsertBlock(block *model.Block, userID string) error
	// @withTransaction
	DeleteBlock(blockID string, modifiedBy string) error
	InsertBlocks(blocks []model.Block, userID string) error
	// @withTransaction
	GetBlockCountsByType() (map[string]int64, error)
	GetBlock(blockID string) (*model.Block, error)
	// @withTransaction
	PatchBlock(blockID string, blockPatch *model.BlockPatch, userID string) error
	GetBlockHistory(blockID string, opts model.QueryBlockHistoryOptions) ([]model.Block, error)
	GetBoardAndCardByID(blockID string) (board *model.Board, card *model.Block, err error)
	GetBoardAndCard(block *model.Block) (board *model.Board, card *model.Block, err error)
	// @withTransaction
	DuplicateBoard(boardID string, userID string, asTemplate bool) (*model.BoardsAndBlocks, []*model.BoardMember, error)
	// @withTransaction
	PatchBlocks(blockPatches *model.BlockPatchBatch, userID string) error

	Shutdown() error

	GetSystemSetting(key string) (string, error)
	GetSystemSettings() (map[string]string, error)
	SetSystemSetting(key, value string) error

	GetRegisteredUserCount() (int, error)
	GetUserByID(userID string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	UpdateUserPassword(username, password string) error
	UpdateUserPasswordByID(userID, password string) error
	GetUsersByTeam(teamID string) ([]*model.User, error)

	GetActiveUserCount(updatedSecondsAgo int64) (int, error)
	GetSession(token string, expireTime int64) (*model.Session, error)
	CreateSession(session *model.Session) error
	RefreshSession(session *model.Session) error
	UpdateSession(session *model.Session) error
	DeleteSession(sessionID string) error
	CleanUpSessions(expireTime int64) error

	UpsertSharing(sharing model.Sharing) error
	GetSharing(rootID string) (*model.Sharing, error)

	UpsertTeamSignupToken(team model.Team) error
	UpsertTeamSettings(team model.Team) error
	GetTeam(ID string) (*model.Team, error)
	GetTeamsForUser(userID string) ([]*model.Team, error)
	GetAllTeams() ([]*model.Team, error)
	GetTeamCount() (int64, error)

	InsertBoard(board *model.Board, userID string) (*model.Board, error)
	// @withTransaction
	InsertBoardWithAdmin(board *model.Board, userID string) (*model.Board, *model.BoardMember, error)
	// @withTransaction
	PatchBoard(boardID string, boardPatch *model.BoardPatch, userID string) (*model.Board, error)
	GetBoard(id string) (*model.Board, error)
	GetBoardsForUserAndTeam(userID, teamID string) ([]*model.Board, error)
	// @withTransaction
	DeleteBoard(boardID, userID string) error

	SaveMember(bm *model.BoardMember) (*model.BoardMember, error)
	DeleteMember(boardID, userID string) error
	GetMemberForBoard(boardID, userID string) (*model.BoardMember, error)
	GetMembersForBoard(boardID string) ([]*model.BoardMember, error)
	SearchBoardsForUserAndTeam(term, userID, teamID string) ([]*model.Board, error)

	// @withTransaction
	CreateBoardsAndBlocksWithAdmin(bab *model.BoardsAndBlocks, userID string) (*model.BoardsAndBlocks, []*model.BoardMember, error)
	// @withTransaction
	CreateBoardsAndBlocks(bab *model.BoardsAndBlocks, userID string) (*model.BoardsAndBlocks, error)
	// @withTransaction
	PatchBoardsAndBlocks(pbab *model.PatchBoardsAndBlocks, userID string) (*model.BoardsAndBlocks, error)
	// @withTransaction
	DeleteBoardsAndBlocks(dbab *model.DeleteBoardsAndBlocks, userID string) error

	GetCategory(id string) (*model.Category, error)
	CreateCategory(category model.Category) error
	UpdateCategory(category model.Category) error
	DeleteCategory(categoryID, userID, teamID string) error

	GetUserCategoryBlocks(userID, teamID string) ([]model.CategoryBlocks, error)
	AddUpdateCategoryBlock(userID, categoryID, blockID string) error

	CreateSubscription(sub *model.Subscription) (*model.Subscription, error)
	DeleteSubscription(blockID string, subscriberID string) error
	GetSubscription(blockID string, subscriberID string) (*model.Subscription, error)
	GetSubscriptions(subscriberID string) ([]*model.Subscription, error)
	GetSubscribersForBlock(blockID string) ([]*model.Subscriber, error)
	GetSubscribersCountForBlock(blockID string) (int, error)
	UpdateSubscribersNotifiedAt(blockID string, notifiedAt int64) error

	UpsertNotificationHint(hint *model.NotificationHint, notificationFreq time.Duration) (*model.NotificationHint, error)
	DeleteNotificationHint(blockID string) error
	GetNotificationHint(blockID string) (*model.NotificationHint, error)
	GetNextNotificationHint(remove bool) (*model.NotificationHint, error)

	RemoveDefaultTemplates(boards []*model.Board) error
	GetTemplateBoards(teamID string) ([]*model.Board, error)

	IsErrNotFound(err error) bool
}

// ErrNotFound is an error type that can be returned by store APIs when a query unexpectedly fetches no records.
type ErrNotFound struct {
	resource string
}

// NewErrNotFound creates a new ErrNotFound instance.
func NewErrNotFound(resource string) *ErrNotFound {
	return &ErrNotFound{
		resource: resource,
	}
}

func (nf *ErrNotFound) Error() string {
	return fmt.Sprintf("{%s} not found", nf.resource)
}

// IsErrNotFound returns true if `err` is or wraps a ErrNotFound.
func IsErrNotFound(err error) bool {
	if err == nil {
		return false
	}

	var nf *ErrNotFound
	return errors.As(err, &nf)
}
