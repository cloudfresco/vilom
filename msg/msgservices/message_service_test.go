package msgservices

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/cloudfresco/vilom/testhelpers"
)

func TestMessageService_GetMessage(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	timeat, err := time.Parse(Layout, "2019-07-23T10:04:26Z")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	messageService := NewMessageService(dbService, redisService)
	msg := Message{}
	msg.ID = uint(1)
	msg.UUID4 = []byte{137, 25, 62, 199, 70, 158, 69, 128, 139, 206, 230, 140, 235, 90, 162, 1}
	msg.IDS = "89193ec7-469e-4580-8bce-e68ceb5aa201"
	msg.NumLikes = uint(0)
	msg.NumUpvotes = uint(0)
	msg.NumDownvotes = uint(0)
	msg.CategoryID = uint(2)
	msg.TopicID = uint(1)
	msg.UserID = uint(1)
	msg.UgroupID = uint(0)
	msg.Statusc = uint(1)
	msg.CreatedAt = timeat
	msg.UpdatedAt = timeat
	msg.CreatedDay = uint(204)
	msg.CreatedWeek = uint(30)
	msg.CreatedMonth = uint(7)
	msg.CreatedYear = uint(2019)
	msg.UpdatedDay = uint(204)
	msg.UpdatedWeek = uint(30)
	msg.UpdatedMonth = uint(7)
	msg.UpdatedYear = uint(2019)

	msgtxt := MessageText{}
	msgtxt.ID = uint(1)
	msgtxt.UUID4 = []byte{148, 157, 162, 242, 107, 10, 67, 245, 165, 223, 218, 71, 208, 232, 206, 72}
	msgtxt.Mtext = "Hi. I am looking into buying a Floptical Drive, and was wondering what experience people have with the drives from Iomega, PLI, MASS MicroSystems, or Procom. These seem to be the main drives on the market. Any advice? Also, I heard about some article in MacWorld about Flopticals. Could someone post a summary, if they have it? Thanks in advance"
	msgtxt.CategoryID = uint(2)
	msgtxt.TopicID = uint(1)
	msgtxt.MessageID = uint(1)
	msgtxt.UserID = uint(1)
	msgtxt.UgroupID = uint(0)
	msgtxt.Statusc = uint(1)
	msgtxt.CreatedAt = timeat
	msgtxt.UpdatedAt = timeat
	msgtxt.CreatedDay = uint(204)
	msgtxt.CreatedWeek = uint(30)
	msgtxt.CreatedMonth = uint(7)
	msgtxt.CreatedYear = uint(2019)
	msgtxt.UpdatedDay = uint(204)
	msgtxt.UpdatedWeek = uint(30)
	msgtxt.UpdatedMonth = uint(7)
	msgtxt.UpdatedYear = uint(2019)

	msg.MessageTexts = append(msg.MessageTexts, &msgtxt)

	msgath := MessageAttachment{}
	msgath.ID = uint(1)
	msgath.UUID4 = []byte{168, 198, 217, 152, 220, 39, 77, 46, 183, 38, 82, 96, 55, 91, 140, 135}
	msgath.Mattach = "mattach"
	msgath.CategoryID = uint(2)
	msgath.TopicID = uint(1)
	msgath.MessageID = uint(1)
	msgath.UserID = uint(1)
	msgath.UgroupID = uint(0)
	msgath.Statusc = uint(1)
	msgath.CreatedAt = timeat
	msgath.UpdatedAt = timeat
	msgath.CreatedDay = uint(204)
	msgath.CreatedWeek = uint(30)
	msgath.CreatedMonth = uint(7)
	msgath.CreatedYear = uint(2019)
	msgath.UpdatedDay = uint(204)
	msgath.UpdatedWeek = uint(30)
	msgath.UpdatedMonth = uint(7)
	msgath.UpdatedYear = uint(2019)

	msg.MessageAttachments = append(msg.MessageAttachments, &msgath)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *MessageService
		args    args
		want    *Message
		wantErr bool
	}{
		{
			t: messageService,
			args: args{
				ctx:       ctx,
				ID:        "89193ec7-469e-4580-8bce-e68ceb5aa201",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &msg,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetMessage(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("MessageService.GetMessage() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("MessageService.GetMessage() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMessageService_GetMessagesWithTextAttach(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	timeat, err := time.Parse(Layout, "2019-07-23T10:04:26Z")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	messageService := NewMessageService(dbService, redisService)
	messages := []*Message{}
	opMessages := []*Message{}
	msg := Message{}
	msg.ID = uint(1)
	msg.UUID4 = []byte{137, 25, 62, 199, 70, 158, 69, 128, 139, 206, 230, 140, 235, 90, 162, 1}
	msg.IDS = "89193ec7-469e-4580-8bce-e68ceb5aa201"
	msg.NumLikes = uint(0)
	msg.NumUpvotes = uint(0)
	msg.NumDownvotes = uint(0)
	msg.CategoryID = uint(2)
	msg.TopicID = uint(1)
	msg.UserID = uint(1)
	msg.UgroupID = uint(0)
	msg.Statusc = uint(1)
	msg.CreatedAt = timeat
	msg.UpdatedAt = timeat
	msg.CreatedDay = uint(204)
	msg.CreatedWeek = uint(30)
	msg.CreatedMonth = uint(7)
	msg.CreatedYear = uint(2019)
	msg.UpdatedDay = uint(204)
	msg.UpdatedWeek = uint(30)
	msg.UpdatedMonth = uint(7)
	msg.UpdatedYear = uint(2019)
	messages = append(messages, &msg)

	opMsg := msg
	msgtxt := MessageText{}
	msgtxt.ID = uint(1)
	msgtxt.UUID4 = []byte{148, 157, 162, 242, 107, 10, 67, 245, 165, 223, 218, 71, 208, 232, 206, 72}
	msgtxt.Mtext = "Hi. I am looking into buying a Floptical Drive, and was wondering what experience people have with the drives from Iomega, PLI, MASS MicroSystems, or Procom. These seem to be the main drives on the market. Any advice? Also, I heard about some article in MacWorld about Flopticals. Could someone post a summary, if they have it? Thanks in advance"
	msgtxt.CategoryID = uint(2)
	msgtxt.TopicID = uint(1)
	msgtxt.MessageID = uint(1)
	msgtxt.UserID = uint(1)
	msgtxt.UgroupID = uint(0)
	msgtxt.Statusc = uint(1)
	msgtxt.CreatedAt = timeat
	msgtxt.UpdatedAt = timeat
	msgtxt.CreatedDay = uint(204)
	msgtxt.CreatedWeek = uint(30)
	msgtxt.CreatedMonth = uint(7)
	msgtxt.CreatedYear = uint(2019)
	msgtxt.UpdatedDay = uint(204)
	msgtxt.UpdatedWeek = uint(30)
	msgtxt.UpdatedMonth = uint(7)
	msgtxt.UpdatedYear = uint(2019)

	opMsg.MessageTexts = append(opMsg.MessageTexts, &msgtxt)

	msgath := MessageAttachment{}
	msgath.ID = uint(1)
	msgath.UUID4 = []byte{168, 198, 217, 152, 220, 39, 77, 46, 183, 38, 82, 96, 55, 91, 140, 135}
	msgath.Mattach = "mattach"
	msgath.CategoryID = uint(2)
	msgath.TopicID = uint(1)
	msgath.MessageID = uint(1)
	msgath.UserID = uint(1)
	msgath.UgroupID = uint(0)
	msgath.Statusc = uint(1)
	msgath.CreatedAt = timeat
	msgath.UpdatedAt = timeat
	msgath.CreatedDay = uint(204)
	msgath.CreatedWeek = uint(30)
	msgath.CreatedMonth = uint(7)
	msgath.CreatedYear = uint(2019)
	msgath.UpdatedDay = uint(204)
	msgath.UpdatedWeek = uint(30)
	msgath.UpdatedMonth = uint(7)
	msgath.UpdatedYear = uint(2019)

	opMsg.MessageAttachments = append(opMsg.MessageAttachments, &msgath)
	opMessages = append(opMessages, &opMsg)

	type args struct {
		ctx       context.Context
		messages  []*Message
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *MessageService
		args    args
		want    []*Message
		wantErr bool
	}{
		{
			t: messageService,
			args: args{
				ctx:       ctx,
				messages:  messages,
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    opMessages,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		got, err := tt.t.GetMessagesWithTextAttach(tt.args.ctx, tt.args.messages, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("MessageService.GetMessagesWithTextAttach() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("MessageService.GetMessagesWithTextAttach() = %v, want %v", got, tt.want)
		}

	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMessageService_GetMessagesTexts(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	timeat, err := time.Parse(Layout, "2019-07-23T10:04:26Z")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	messageService := NewMessageService(dbService, redisService)

	msgtxt := MessageText{}
	msgtxts := []*MessageText{}
	msgtxt.ID = uint(1)
	msgtxt.UUID4 = []byte{148, 157, 162, 242, 107, 10, 67, 245, 165, 223, 218, 71, 208, 232, 206, 72}
	msgtxt.Mtext = "Hi. I am looking into buying a Floptical Drive, and was wondering what experience people have with the drives from Iomega, PLI, MASS MicroSystems, or Procom. These seem to be the main drives on the market. Any advice? Also, I heard about some article in MacWorld about Flopticals. Could someone post a summary, if they have it? Thanks in advance"
	msgtxt.CategoryID = uint(2)
	msgtxt.TopicID = uint(1)
	msgtxt.MessageID = uint(1)
	msgtxt.UserID = uint(1)
	msgtxt.UgroupID = uint(0)
	msgtxt.Statusc = uint(1)
	msgtxt.CreatedAt = timeat
	msgtxt.UpdatedAt = timeat
	msgtxt.CreatedDay = uint(204)
	msgtxt.CreatedWeek = uint(30)
	msgtxt.CreatedMonth = uint(7)
	msgtxt.CreatedYear = uint(2019)
	msgtxt.UpdatedDay = uint(204)
	msgtxt.UpdatedWeek = uint(30)
	msgtxt.UpdatedMonth = uint(7)
	msgtxt.UpdatedYear = uint(2019)

	msgtxts = append(msgtxts, &msgtxt)

	type args struct {
		ctx       context.Context
		messageID uint
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *MessageService
		args    args
		want    []*MessageText
		wantErr bool
	}{
		{
			t: messageService,
			args: args{
				ctx:       ctx,
				messageID: uint(1),
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    msgtxts,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		got, err := tt.t.GetMessagesTexts(tt.args.ctx, tt.args.messageID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("MessageService.GetMessagesTexts() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("MessageService.GetMessagesTexts() = %v, want %v", got, tt.want)
		}

	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMessageService_GetMessageAttachments(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	timeat, err := time.Parse(Layout, "2019-07-23T10:04:26Z")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	messageService := NewMessageService(dbService, redisService)

	msgaths := []*MessageAttachment{}
	msgath := MessageAttachment{}
	msgath.ID = uint(1)
	msgath.UUID4 = []byte{168, 198, 217, 152, 220, 39, 77, 46, 183, 38, 82, 96, 55, 91, 140, 135}
	msgath.Mattach = "mattach"
	msgath.CategoryID = uint(2)
	msgath.TopicID = uint(1)
	msgath.MessageID = uint(1)
	msgath.UserID = uint(1)
	msgath.UgroupID = uint(0)
	msgath.Statusc = uint(1)
	msgath.CreatedAt = timeat
	msgath.UpdatedAt = timeat
	msgath.CreatedDay = uint(204)
	msgath.CreatedWeek = uint(30)
	msgath.CreatedMonth = uint(7)
	msgath.CreatedYear = uint(2019)
	msgath.UpdatedDay = uint(204)
	msgath.UpdatedWeek = uint(30)
	msgath.UpdatedMonth = uint(7)
	msgath.UpdatedYear = uint(2019)

	msgaths = append(msgaths, &msgath)

	type args struct {
		ctx       context.Context
		messageID uint
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *MessageService
		args    args
		want    []*MessageAttachment
		wantErr bool
	}{
		{
			t: messageService,
			args: args{
				ctx:       ctx,
				messageID: uint(1),
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    msgaths,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		got, err := tt.t.GetMessageAttachments(tt.args.ctx, tt.args.messageID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("MessageService.GetMessageAttachments() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("MessageService.GetMessageAttachments() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMessageService_UpdateMessage(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	messageService := NewMessageService(dbService, redisService)

	form1 := Message{}
	form1.Mtext = "Messagetext2"

	type args struct {
		ctx       context.Context
		ID        string
		form      *Message
		UserID    string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *MessageService
		args    args
		wantErr bool
	}{
		{
			t: messageService,
			args: args{
				ctx:       ctx,
				ID:        "89193ec7-469e-4580-8bce-e68ceb5aa201",
				form:      &form1,
				UserID:    "29ea215b-8fb3-4453-b413-81a661e44495",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := tt.t.UpdateMessage(tt.args.ctx, tt.args.ID, tt.args.form, tt.args.UserID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("MessageService.UpdateMessage() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMessageService_DeleteMessage(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	messageService := NewMessageService(dbService, redisService)
	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *MessageService
		args    args
		wantErr bool
	}{
		{
			t: messageService,
			args: args{
				ctx:       ctx,
				ID:        "89193ec7-469e-4580-8bce-e68ceb5aa201",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := tt.t.DeleteMessage(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("MessageService.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}
