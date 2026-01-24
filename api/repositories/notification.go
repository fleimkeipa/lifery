package repositories

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/pkg/logger"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type NotificationRepository struct {
	db *pg.DB
}

func NewNotificationRepository(db *pg.DB) *NotificationRepository {
	rc := &NotificationRepository{
		db: db,
	}

	if err := rc.createSchema(db); err != nil {
		logger.Log.Fatalf("failed to create schema: %v", err)
	}

	return rc
}

func (rc *NotificationRepository) Create(ctx context.Context, newNotification *model.Notification) (*model.Notification, error) {
	sqlNotification := rc.internalToSQL(newNotification)

	q := rc.db.Model(sqlNotification)

	_, err := q.Insert()
	if err != nil {
		return nil, pkg.NewError(err, "failed to create notification", http.StatusInternalServerError)
	}

	return rc.sqlToInternal(sqlNotification), nil
}

func (rc *NotificationRepository) Update(ctx context.Context, notificationID string, newNotification *model.Notification) (*model.Notification, error) {
	if notificationID == "" || notificationID == "0" {
		return nil, pkg.NewError(nil, "invalid notification id "+notificationID, http.StatusBadRequest)
	}

	newNotification.ID = notificationID

	sqlNotification := rc.internalToSQL(newNotification)

	q := rc.db.Model(sqlNotification).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, pkg.NewError(err, "failed to update notification", http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return nil, pkg.NewError(nil, "no notification updated", http.StatusBadRequest)
	}

	return rc.sqlToInternal(sqlNotification), nil
}

func (rc *NotificationRepository) GetByID(ctx context.Context, notificationID string) (*model.Notification, error) {
	id, _ := strconv.Atoi(notificationID)

	internalNotification := &notification{
		ID: id,
	}

	query := rc.db.Model(internalNotification).Where("id = ?", id)
	if err := query.Select(); err != nil {
		return nil, pkg.NewError(err, "failed to get notification", http.StatusInternalServerError)
	}

	return rc.sqlToInternal(internalNotification), nil
}

func (rc *NotificationRepository) List(ctx context.Context, opts *model.NotificationFindOpts) (*model.NotificationList, error) {
	if opts == nil {
		return nil, pkg.NewError(nil, "opts is nil", http.StatusBadRequest)
	}

	notifications := make([]notification, 0)

	query := rc.db.Model(&notifications)

	query = applyOrderBy(query, opts.OrderByOpts)
	query = applyStandardQueries(query, opts.PaginationOpts)
	query = rc.fillFields(query, opts)
	query = rc.fillNotificationFilter(query, opts)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, pkg.NewError(err, "failed to list notifications", http.StatusInternalServerError)
	}

	internalNotifications := make([]model.Notification, 0)
	for _, v := range notifications {
		internalNotifications = append(internalNotifications, *rc.sqlToInternal(&v))
	}

	return &model.NotificationList{
		Notifications: internalNotifications,
		Total:         count,
		PaginationOpts: model.PaginationOpts{
			Skip:  opts.Skip,
			Limit: opts.Limit,
		},
	}, nil
}

func (rc *NotificationRepository) Delete(ctx context.Context, notificationID string) error {
	id, _ := strconv.Atoi(notificationID)

	internalNotification := &notification{
		ID: id,
	}

	result, err := rc.db.Model(internalNotification).Where("id = ?", id).Delete()
	if err != nil {
		return pkg.NewError(err, "failed to delete notification", http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return pkg.NewError(nil, "no notification deleted", http.StatusBadRequest)
	}

	return nil
}

func (rc *NotificationRepository) fillFields(tx *orm.Query, opts *model.NotificationFindOpts) *orm.Query {
	fields := opts.Fields

	if len(fields) == 0 {
		return tx
	}

	if len(fields) == 1 && fields[0] == model.ZeroCreds {
		return tx.Column(
			"notification.id",
			"notification.user_id",
			"notification.type",
			"notification.message",
			"notification.read",
			"notification.created_at",
		)
	}

	qualifiedFields := make([]string, len(fields))
	for i, field := range fields {
		qualifiedFields[i] = "notification." + field
	}

	return tx.Column(qualifiedFields...)
}

func (rc *NotificationRepository) fillNotificationFilter(tx *orm.Query, opts *model.NotificationFindOpts) *orm.Query {
	if opts.UserID.IsSended {
		tx = applyFilterWithOperand(tx, "user_id", opts.UserID)
	}

	if opts.Read.IsSended {
		tx = applyFilterWithOperand(tx, "read", opts.Read)
	}

	return tx
}

func (rc *NotificationRepository) internalToSQL(newNotification *model.Notification) *notification {
	nID, _ := strconv.Atoi(newNotification.ID)
	userID, _ := strconv.Atoi(newNotification.UserID)
	return &notification{
		Type:      newNotification.Type,
		Message:   newNotification.Message,
		Read:      int(newNotification.Read),
		UserID:    userID,
		ID:        nID,
		CreatedAt: newNotification.CreatedAt.Format(time.RFC3339),
	}
}

func (rc *NotificationRepository) sqlToInternal(notification *notification) *model.Notification {
	createdAt, _ := time.Parse(time.RFC3339, notification.CreatedAt)
	return &model.Notification{
		ID:        strconv.Itoa(notification.ID),
		UserID:    strconv.Itoa(notification.UserID),
		Type:      notification.Type,
		Message:   notification.Message,
		Read:      model.NotificationStatus(notification.Read),
		CreatedAt: createdAt,
	}
}

func (rc *NotificationRepository) createSchema(db *pg.DB) error {
	model := (*notification)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return pkg.NewError(err, "failed to create notification table", http.StatusInternalServerError)
	}

	return nil
}
