package msgservices

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/cloudfresco/vilom/testhelpers"
)

func TestTopicService_ShowTopic(t *testing.T) {
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
	topicService := NewTopicService(dbService, redisService)
	topic := Topic{}
	topic.ID = uint(1)
	topic.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	topic.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	topic.TopicName = "Floptical Question"
	topic.TopicDesc = "Floptical Question"
	topic.NumTags = 0
	topic.Tag1 = ""
	topic.Tag2 = ""
	topic.Tag3 = ""
	topic.Tag4 = ""
	topic.Tag5 = ""
	topic.Tag6 = ""
	topic.Tag7 = ""
	topic.Tag8 = ""
	topic.Tag9 = ""
	topic.Tag10 = ""
	topic.NumViews = uint(0)
	topic.NumMessages = uint(1)
	topic.CategoryID = uint(2)
	topic.UserID = uint(1)
	topic.UgroupID = uint(0)
	topic.Statusc = uint(1)
	topic.CreatedAt = timeat
	topic.UpdatedAt = timeat
	topic.CreatedDay = uint(204)
	topic.CreatedWeek = uint(30)
	topic.CreatedMonth = uint(7)
	topic.CreatedYear = uint(2019)
	topic.UpdatedDay = uint(204)
	topic.UpdatedWeek = uint(30)
	topic.UpdatedMonth = uint(7)
	topic.UpdatedYear = uint(2019)

	messages := []*Message{}
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
	messages = append(messages, &msg)

	topic.Messages = messages

	type args struct {
		ctx       context.Context
		ID        string
		UserID    string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		want    *Topic
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				UserID:    "29ea215b-8fb3-4453-b413-81a661e44495",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &topic,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.ShowTopic(tt.args.ctx, tt.args.ID, tt.args.UserID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TopicService.ShowTopic() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopicService.ShowTopic() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_GetTopicByID(t *testing.T) {
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
	topicService := NewTopicService(dbService, redisService)
	topic := Topic{}
	topic.ID = uint(1)
	topic.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	topic.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	topic.TopicName = "Floptical Question"
	topic.TopicDesc = "Floptical Question"
	topic.NumTags = 0
	topic.Tag1 = ""
	topic.Tag2 = ""
	topic.Tag3 = ""
	topic.Tag4 = ""
	topic.Tag5 = ""
	topic.Tag6 = ""
	topic.Tag7 = ""
	topic.Tag8 = ""
	topic.Tag9 = ""
	topic.Tag10 = ""
	topic.NumViews = uint(0)
	topic.NumMessages = uint(1)
	topic.CategoryID = uint(2)
	topic.UserID = uint(1)
	topic.UgroupID = uint(0)
	topic.Statusc = uint(1)
	topic.CreatedAt = timeat
	topic.UpdatedAt = timeat
	topic.CreatedDay = uint(204)
	topic.CreatedWeek = uint(30)
	topic.CreatedMonth = uint(7)
	topic.CreatedYear = uint(2019)
	topic.UpdatedDay = uint(204)
	topic.UpdatedWeek = uint(30)
	topic.UpdatedMonth = uint(7)
	topic.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        uint
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		want    *Topic
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				ID:        uint(1),
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &topic,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetTopicByID(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TopicService.GetTopicByID() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopicService.GetTopicByID() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_GetTopic(t *testing.T) {
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
	topicService := NewTopicService(dbService, redisService)
	topic := Topic{}
	topic.ID = uint(1)
	topic.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	topic.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	topic.TopicName = "Floptical Question"
	topic.TopicDesc = "Floptical Question"
	topic.NumTags = 0
	topic.Tag1 = ""
	topic.Tag2 = ""
	topic.Tag3 = ""
	topic.Tag4 = ""
	topic.Tag5 = ""
	topic.Tag6 = ""
	topic.Tag7 = ""
	topic.Tag8 = ""
	topic.Tag9 = ""
	topic.Tag10 = ""
	topic.NumViews = uint(0)
	topic.NumMessages = uint(1)
	topic.CategoryID = uint(2)
	topic.UserID = uint(1)
	topic.UgroupID = uint(0)
	topic.Statusc = uint(1)
	topic.CreatedAt = timeat
	topic.UpdatedAt = timeat
	topic.CreatedDay = uint(204)
	topic.CreatedWeek = uint(30)
	topic.CreatedMonth = uint(7)
	topic.CreatedYear = uint(2019)
	topic.UpdatedDay = uint(204)
	topic.UpdatedWeek = uint(30)
	topic.UpdatedMonth = uint(7)
	topic.UpdatedYear = uint(2019)
	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		want    *Topic
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &topic,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		got, err := tt.t.GetTopic(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TopicService.GetTopic() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopicService.GetTopic() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_GetTopicByName(t *testing.T) {
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
	topicService := NewTopicService(dbService, redisService)
	topic := Topic{}
	topic.ID = uint(1)
	topic.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	topic.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	topic.TopicName = "Floptical Question"
	topic.TopicDesc = "Floptical Question"
	topic.NumTags = 0
	topic.Tag1 = ""
	topic.Tag2 = ""
	topic.Tag3 = ""
	topic.Tag4 = ""
	topic.Tag5 = ""
	topic.Tag6 = ""
	topic.Tag7 = ""
	topic.Tag8 = ""
	topic.Tag9 = ""
	topic.Tag10 = ""
	topic.NumViews = uint(0)
	topic.NumMessages = uint(1)
	topic.CategoryID = uint(2)
	topic.UserID = uint(1)
	topic.UgroupID = uint(0)
	topic.Statusc = uint(1)
	topic.CreatedAt = timeat
	topic.UpdatedAt = timeat
	topic.CreatedDay = uint(204)
	topic.CreatedWeek = uint(30)
	topic.CreatedMonth = uint(7)
	topic.CreatedYear = uint(2019)
	topic.UpdatedDay = uint(204)
	topic.UpdatedWeek = uint(30)
	topic.UpdatedMonth = uint(7)
	topic.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		topicname string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		want    *Topic
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				topicname: "Floptical Question",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &topic,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetTopicByName(tt.args.ctx, tt.args.topicname, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TopicService.GetTopicByName() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopicService.GetTopicByName() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_GetTopicWithMessages(t *testing.T) {
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
	topicService := NewTopicService(dbService, redisService)
	topic := Topic{}
	topic.ID = uint(1)
	topic.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	topic.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	topic.TopicName = "Floptical Question"
	topic.TopicDesc = "Floptical Question"
	topic.NumTags = 0
	topic.Tag1 = ""
	topic.Tag2 = ""
	topic.Tag3 = ""
	topic.Tag4 = ""
	topic.Tag5 = ""
	topic.Tag6 = ""
	topic.Tag7 = ""
	topic.Tag8 = ""
	topic.Tag9 = ""
	topic.Tag10 = ""
	topic.NumViews = uint(0)
	topic.NumMessages = uint(1)
	topic.CategoryID = uint(2)
	topic.UserID = uint(1)
	topic.UgroupID = uint(0)
	topic.Statusc = uint(1)
	topic.CreatedAt = timeat
	topic.UpdatedAt = timeat
	topic.CreatedDay = uint(204)
	topic.CreatedWeek = uint(30)
	topic.CreatedMonth = uint(7)
	topic.CreatedYear = uint(2019)
	topic.UpdatedDay = uint(204)
	topic.UpdatedWeek = uint(30)
	topic.UpdatedMonth = uint(7)
	topic.UpdatedYear = uint(2019)
	messages := []*Message{}
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
	messages = append(messages, &msg)

	topic.Messages = messages

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		want    *Topic
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &topic,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetTopicWithMessages(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TopicService.GetTopicWithMessages() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopicService.GetTopicWithMessages() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_GetTopicMessages(t *testing.T) {
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
	topicService := NewTopicService(dbService, redisService)
	topic := Topic{}
	topic.ID = uint(1)
	topic.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	topic.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	topic.TopicName = "Floptical Question"
	topic.TopicDesc = "Floptical Question"
	topic.NumTags = 0
	topic.Tag1 = ""
	topic.Tag2 = ""
	topic.Tag3 = ""
	topic.Tag4 = ""
	topic.Tag5 = ""
	topic.Tag6 = ""
	topic.Tag7 = ""
	topic.Tag8 = ""
	topic.Tag9 = ""
	topic.Tag10 = ""
	topic.NumViews = uint(0)
	topic.NumMessages = uint(1)
	topic.CategoryID = uint(2)
	topic.UserID = uint(1)
	topic.UgroupID = uint(0)
	topic.Statusc = uint(1)
	topic.CreatedAt = timeat
	topic.UpdatedAt = timeat
	topic.CreatedDay = uint(204)
	topic.CreatedWeek = uint(30)
	topic.CreatedMonth = uint(7)
	topic.CreatedYear = uint(2019)
	topic.UpdatedDay = uint(204)
	topic.UpdatedWeek = uint(30)
	topic.UpdatedMonth = uint(7)
	topic.UpdatedYear = uint(2019)

	messages := []*Message{}
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

	topic.Messages = messages

	type args struct {
		ctx       context.Context
		uuid4byte []byte
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		want    *Topic
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				uuid4byte: []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195},
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &topic,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetTopicMessages(tt.args.ctx, tt.args.uuid4byte, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TopicService.GetTopicMessages() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopicService.GetTopicMessages() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_GetTopicsUser(t *testing.T) {
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
	topicService := NewTopicService(dbService, redisService)

	tu := TopicsUser{}
	tu.ID = uint(1)
	tu.UUID4 = []byte{13, 46, 69, 184, 91, 15, 78, 37, 172, 121, 102, 33, 237, 119, 10, 77}
	tu.IDS = "0d2e45b8-5b0f-4e25-ac79-6621ed770a4d"
	tu.TopicID = uint(1)
	tu.NumMessages = uint(0)
	tu.NumViews = 1
	tu.UserID = uint(1)
	tu.UgroupID = uint(0)
	tu.Statusc = uint(1)
	tu.CreatedAt = timeat
	tu.UpdatedAt = timeat
	tu.CreatedDay = uint(204)
	tu.CreatedWeek = uint(30)
	tu.CreatedMonth = uint(7)
	tu.CreatedYear = uint(2019)
	tu.UpdatedDay = uint(204)
	tu.UpdatedWeek = uint(30)
	tu.UpdatedMonth = uint(7)
	tu.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        uint
		UserID    uint
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		want    *TopicsUser
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				ID:        uint(1),
				UserID:    uint(1),
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &tu,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetTopicsUser(tt.args.ctx, tt.args.ID, tt.args.UserID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TopicService.GetTopicsUser() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopicService.GetTopicsUser() = %v, want %v", got, tt.want)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_UpdateTopic(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	topicService := NewTopicService(dbService, redisService)
	form1 := Topic{}
	form1.TopicName = "Hard drive security2"
	form1.TopicDesc = "Hard drive security2"
	type args struct {
		ctx       context.Context
		ID        string
		form      *Topic
		UserID    string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				form:      &form1,
				UserID:    "29ea215b-8fb3-4453-b413-81a661e44495",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := tt.t.UpdateTopic(tt.args.ctx, tt.args.ID, tt.args.form, tt.args.UserID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("TopicService.UpdateTopic() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTopicService_DeleteTopic(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	topicService := NewTopicService(dbService, redisService)
	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *TopicService
		args    args
		wantErr bool
	}{
		{
			t: topicService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := tt.t.DeleteTopic(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("TopicService.DeleteTopic() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
	err = testhelpers.DeleteSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}
}
