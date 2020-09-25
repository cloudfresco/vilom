package msgservices

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/cloudfresco/vilom/testhelpers"
)

func TestChannelService_ShowChannel(t *testing.T) {
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
	channelService := NewChannelService(dbService, redisService)
	channel := Channel{}
	channel.ID = uint(1)
	channel.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	channel.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	channel.ChannelName = "Floptical Question"
	channel.ChannelDesc = "Floptical Question"
	channel.NumTags = 0
	channel.Tag1 = ""
	channel.Tag2 = ""
	channel.Tag3 = ""
	channel.Tag4 = ""
	channel.Tag5 = ""
	channel.Tag6 = ""
	channel.Tag7 = ""
	channel.Tag8 = ""
	channel.Tag9 = ""
	channel.Tag10 = ""
	channel.NumViews = uint(0)
	channel.NumMessages = uint(1)
	channel.WorkspaceID = uint(2)
	channel.UserID = uint(1)
	channel.UgroupID = uint(0)
	channel.Statusc = uint(1)
	channel.CreatedAt = timeat
	channel.UpdatedAt = timeat
	channel.CreatedDay = uint(204)
	channel.CreatedWeek = uint(30)
	channel.CreatedMonth = uint(7)
	channel.CreatedYear = uint(2019)
	channel.UpdatedDay = uint(204)
	channel.UpdatedWeek = uint(30)
	channel.UpdatedMonth = uint(7)
	channel.UpdatedYear = uint(2019)

	messages := []*Message{}
	msg := Message{}
	msg.ID = uint(1)
	msg.UUID4 = []byte{137, 25, 62, 199, 70, 158, 69, 128, 139, 206, 230, 140, 235, 90, 162, 1}
	msg.IDS = "89193ec7-469e-4580-8bce-e68ceb5aa201"
	msg.NumLikes = uint(0)
	msg.NumUpvotes = uint(0)
	msg.NumDownvotes = uint(0)
	msg.WorkspaceID = uint(2)
	msg.ChannelID = uint(1)
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
	msgtxt.WorkspaceID = uint(2)
	msgtxt.ChannelID = uint(1)
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
	msgath.WorkspaceID = uint(2)
	msgath.ChannelID = uint(1)
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

	channel.Messages = messages

	type args struct {
		ctx       context.Context
		ID        string
		UserID    string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		want    *Channel
		wantErr bool
	}{
		{
			t: channelService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				UserID:    "29ea215b-8fb3-4453-b413-81a661e44495",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &channel,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.ShowChannel(tt.args.ctx, tt.args.ID, tt.args.UserID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.ShowChannel() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChannelService.ShowChannel() = %v, want %v", got, tt.want)
		}
	}
}

func TestChannelService_GetChannelByID(t *testing.T) {
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
	channelService := NewChannelService(dbService, redisService)
	channel := Channel{}
	channel.ID = uint(1)
	channel.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	channel.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	channel.ChannelName = "Floptical Question"
	channel.ChannelDesc = "Floptical Question"
	channel.NumTags = 0
	channel.Tag1 = ""
	channel.Tag2 = ""
	channel.Tag3 = ""
	channel.Tag4 = ""
	channel.Tag5 = ""
	channel.Tag6 = ""
	channel.Tag7 = ""
	channel.Tag8 = ""
	channel.Tag9 = ""
	channel.Tag10 = ""
	channel.NumViews = uint(0)
	channel.NumMessages = uint(1)
	channel.WorkspaceID = uint(2)
	channel.UserID = uint(1)
	channel.UgroupID = uint(0)
	channel.Statusc = uint(1)
	channel.CreatedAt = timeat
	channel.UpdatedAt = timeat
	channel.CreatedDay = uint(204)
	channel.CreatedWeek = uint(30)
	channel.CreatedMonth = uint(7)
	channel.CreatedYear = uint(2019)
	channel.UpdatedDay = uint(204)
	channel.UpdatedWeek = uint(30)
	channel.UpdatedMonth = uint(7)
	channel.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        uint
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		want    *Channel
		wantErr bool
	}{
		{
			t: channelService,
			args: args{
				ctx:       ctx,
				ID:        uint(1),
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &channel,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetChannelByID(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.GetChannelByID() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChannelService.GetChannelByID() = %v, want %v", got, tt.want)
		}
	}
}

func TestChannelService_GetChannel(t *testing.T) {
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
	channelService := NewChannelService(dbService, redisService)
	channel := Channel{}
	channel.ID = uint(1)
	channel.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	channel.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	channel.ChannelName = "Floptical Question"
	channel.ChannelDesc = "Floptical Question"
	channel.NumTags = 0
	channel.Tag1 = ""
	channel.Tag2 = ""
	channel.Tag3 = ""
	channel.Tag4 = ""
	channel.Tag5 = ""
	channel.Tag6 = ""
	channel.Tag7 = ""
	channel.Tag8 = ""
	channel.Tag9 = ""
	channel.Tag10 = ""
	channel.NumViews = uint(0)
	channel.NumMessages = uint(1)
	channel.WorkspaceID = uint(2)
	channel.UserID = uint(1)
	channel.UgroupID = uint(0)
	channel.Statusc = uint(1)
	channel.CreatedAt = timeat
	channel.UpdatedAt = timeat
	channel.CreatedDay = uint(204)
	channel.CreatedWeek = uint(30)
	channel.CreatedMonth = uint(7)
	channel.CreatedYear = uint(2019)
	channel.UpdatedDay = uint(204)
	channel.UpdatedWeek = uint(30)
	channel.UpdatedMonth = uint(7)
	channel.UpdatedYear = uint(2019)
	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		want    *Channel
		wantErr bool
	}{
		{
			t: channelService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &channel,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		got, err := tt.t.GetChannel(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.GetChannel() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChannelService.GetChannel() = %v, want %v", got, tt.want)
		}
	}
}

func TestChannelService_GetChannelByName(t *testing.T) {
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
	channelService := NewChannelService(dbService, redisService)
	channel := Channel{}
	channel.ID = uint(1)
	channel.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	channel.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	channel.ChannelName = "Floptical Question"
	channel.ChannelDesc = "Floptical Question"
	channel.NumTags = 0
	channel.Tag1 = ""
	channel.Tag2 = ""
	channel.Tag3 = ""
	channel.Tag4 = ""
	channel.Tag5 = ""
	channel.Tag6 = ""
	channel.Tag7 = ""
	channel.Tag8 = ""
	channel.Tag9 = ""
	channel.Tag10 = ""
	channel.NumViews = uint(0)
	channel.NumMessages = uint(1)
	channel.WorkspaceID = uint(2)
	channel.UserID = uint(1)
	channel.UgroupID = uint(0)
	channel.Statusc = uint(1)
	channel.CreatedAt = timeat
	channel.UpdatedAt = timeat
	channel.CreatedDay = uint(204)
	channel.CreatedWeek = uint(30)
	channel.CreatedMonth = uint(7)
	channel.CreatedYear = uint(2019)
	channel.UpdatedDay = uint(204)
	channel.UpdatedWeek = uint(30)
	channel.UpdatedMonth = uint(7)
	channel.UpdatedYear = uint(2019)

	type args struct {
		ctx         context.Context
		channelname string
		userEmail   string
		requestID   string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		want    *Channel
		wantErr bool
	}{
		{
			t: channelService,
			args: args{
				ctx:         ctx,
				channelname: "Floptical Question",
				userEmail:   "abcd145@gmail.com",
				requestID:   "bks1m1g91jau4nkks2f0",
			},
			want:    &channel,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetChannelByName(tt.args.ctx, tt.args.channelname, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.GetChannelByName() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChannelService.GetChannelByName() = %v, want %v", got, tt.want)
		}
	}
}

func TestChannelService_GetChannelWithMessages(t *testing.T) {
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
	channelService := NewChannelService(dbService, redisService)
	channel := Channel{}
	channel.ID = uint(1)
	channel.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	channel.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	channel.ChannelName = "Floptical Question"
	channel.ChannelDesc = "Floptical Question"
	channel.NumTags = 0
	channel.Tag1 = ""
	channel.Tag2 = ""
	channel.Tag3 = ""
	channel.Tag4 = ""
	channel.Tag5 = ""
	channel.Tag6 = ""
	channel.Tag7 = ""
	channel.Tag8 = ""
	channel.Tag9 = ""
	channel.Tag10 = ""
	channel.NumViews = uint(0)
	channel.NumMessages = uint(1)
	channel.WorkspaceID = uint(2)
	channel.UserID = uint(1)
	channel.UgroupID = uint(0)
	channel.Statusc = uint(1)
	channel.CreatedAt = timeat
	channel.UpdatedAt = timeat
	channel.CreatedDay = uint(204)
	channel.CreatedWeek = uint(30)
	channel.CreatedMonth = uint(7)
	channel.CreatedYear = uint(2019)
	channel.UpdatedDay = uint(204)
	channel.UpdatedWeek = uint(30)
	channel.UpdatedMonth = uint(7)
	channel.UpdatedYear = uint(2019)
	messages := []*Message{}
	msg := Message{}
	msg.ID = uint(1)
	msg.UUID4 = []byte{137, 25, 62, 199, 70, 158, 69, 128, 139, 206, 230, 140, 235, 90, 162, 1}
	msg.IDS = "89193ec7-469e-4580-8bce-e68ceb5aa201"
	msg.NumLikes = uint(0)
	msg.NumUpvotes = uint(0)
	msg.NumDownvotes = uint(0)
	msg.WorkspaceID = uint(2)
	msg.ChannelID = uint(1)
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
	msgtxt.WorkspaceID = uint(2)
	msgtxt.ChannelID = uint(1)
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
	msgath.WorkspaceID = uint(2)
	msgath.ChannelID = uint(1)
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

	channel.Messages = messages

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		want    *Channel
		wantErr bool
	}{
		{
			t: channelService,
			args: args{
				ctx:       ctx,
				ID:        "44b2e674-7031-4487-be96-60093bfe8ac3",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &channel,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetChannelWithMessages(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.GetChannelWithMessages() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChannelService.GetChannelWithMessages() = %v, want %v", got, tt.want)
		}
	}
}

func TestChannelService_GetChannelMessages(t *testing.T) {
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
	channelService := NewChannelService(dbService, redisService)
	channel := Channel{}
	channel.ID = uint(1)
	channel.UUID4 = []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195}
	channel.IDS = "44b2e674-7031-4487-be96-60093bfe8ac3"
	channel.ChannelName = "Floptical Question"
	channel.ChannelDesc = "Floptical Question"
	channel.NumTags = 0
	channel.Tag1 = ""
	channel.Tag2 = ""
	channel.Tag3 = ""
	channel.Tag4 = ""
	channel.Tag5 = ""
	channel.Tag6 = ""
	channel.Tag7 = ""
	channel.Tag8 = ""
	channel.Tag9 = ""
	channel.Tag10 = ""
	channel.NumViews = uint(0)
	channel.NumMessages = uint(1)
	channel.WorkspaceID = uint(2)
	channel.UserID = uint(1)
	channel.UgroupID = uint(0)
	channel.Statusc = uint(1)
	channel.CreatedAt = timeat
	channel.UpdatedAt = timeat
	channel.CreatedDay = uint(204)
	channel.CreatedWeek = uint(30)
	channel.CreatedMonth = uint(7)
	channel.CreatedYear = uint(2019)
	channel.UpdatedDay = uint(204)
	channel.UpdatedWeek = uint(30)
	channel.UpdatedMonth = uint(7)
	channel.UpdatedYear = uint(2019)

	messages := []*Message{}
	msg := Message{}
	msg.ID = uint(1)
	msg.UUID4 = []byte{137, 25, 62, 199, 70, 158, 69, 128, 139, 206, 230, 140, 235, 90, 162, 1}
	msg.IDS = "89193ec7-469e-4580-8bce-e68ceb5aa201"
	msg.NumLikes = uint(0)
	msg.NumUpvotes = uint(0)
	msg.NumDownvotes = uint(0)
	msg.WorkspaceID = uint(2)
	msg.ChannelID = uint(1)
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

	channel.Messages = messages

	type args struct {
		ctx       context.Context
		uuid4byte []byte
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		want    *Channel
		wantErr bool
	}{
		{
			t: channelService,
			args: args{
				ctx:       ctx,
				uuid4byte: []byte{68, 178, 230, 116, 112, 49, 68, 135, 190, 150, 96, 9, 59, 254, 138, 195},
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &channel,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.GetChannelMessages(tt.args.ctx, tt.args.uuid4byte, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.GetChannelMessages() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChannelService.GetChannelMessages() = %v, want %v", got, tt.want)
		}
	}
}

func TestChannelService_GetChannelsUser(t *testing.T) {
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
	channelService := NewChannelService(dbService, redisService)

	tu := ChannelsUser{}
	tu.ID = uint(1)
	tu.UUID4 = []byte{13, 46, 69, 184, 91, 15, 78, 37, 172, 121, 102, 33, 237, 119, 10, 77}
	tu.IDS = "0d2e45b8-5b0f-4e25-ac79-6621ed770a4d"
	tu.ChannelID = uint(1)
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
		t       *ChannelService
		args    args
		want    *ChannelsUser
		wantErr bool
	}{
		{
			t: channelService,
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
		got, err := tt.t.GetChannelsUser(tt.args.ctx, tt.args.ID, tt.args.UserID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.GetChannelsUser() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChannelService.GetChannelsUser() = %v, want %v", got, tt.want)
		}
	}
}

func TestChannelService_UpdateChannel(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	channelService := NewChannelService(dbService, redisService)
	form1 := Channel{}
	form1.ChannelName = "Hard drive security2"
	form1.ChannelDesc = "Hard drive security2"
	type args struct {
		ctx       context.Context
		ID        string
		form      *Channel
		UserID    string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		wantErr bool
	}{
		{
			t: channelService,
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
		if err := tt.t.UpdateChannel(tt.args.ctx, tt.args.ID, tt.args.form, tt.args.UserID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.UpdateChannel() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestChannelService_DeleteChannel(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	channelService := NewChannelService(dbService, redisService)
	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		t       *ChannelService
		args    args
		wantErr bool
	}{
		{
			t: channelService,
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
		if err := tt.t.DeleteChannel(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("ChannelService.DeleteChannel() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
