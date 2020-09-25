package msgservices

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/cloudfresco/vilom/testhelpers"
	_ "github.com/go-sql-driver/mysql"
)

func TestWorkspaceService_GetWorkspaces(t *testing.T) {
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
	workspaceService := NewWorkspaceService(dbService, redisService)
	workspace := Workspace{}
	workspaces := []*Workspace{}
	workspace.ID = uint(1)
	workspace.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	workspace.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	workspace.WorkspaceName = "Performance Portable Transmitter"
	workspace.WorkspaceDesc = "Performance Portable Transmitter"
	workspace.NumViews = uint(0)
	workspace.NumChannels = uint(0)
	workspace.Levelc = uint(0)
	workspace.ParentID = uint(0)
	workspace.NumChd = uint(1)
	workspace.UgroupID = uint(0)
	workspace.UserID = uint(1)
	workspace.Statusc = uint(1)
	workspace.CreatedAt = timeat
	workspace.UpdatedAt = timeat
	workspace.CreatedDay = uint(204)
	workspace.CreatedWeek = uint(30)
	workspace.CreatedMonth = uint(7)
	workspace.CreatedYear = uint(2019)
	workspace.UpdatedDay = uint(204)
	workspace.UpdatedWeek = uint(30)
	workspace.UpdatedMonth = uint(7)
	workspace.UpdatedYear = uint(2019)

	workspaces = append(workspaces, &workspace)
	nextc := "MA=="
	c1 := WorkspaceCursor{workspaces, nextc}

	type args struct {
		ctx        context.Context
		limit      string
		nextCursor string
		userEmail  string
		requestID  string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		want    *WorkspaceCursor
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:        ctx,
				limit:      "4",
				nextCursor: "",
				userEmail:  "abcd145@gmail.com",
				requestID:  "bks1m1g91jau4nkks2f0",
			},
			want:    &c1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetWorkspaces(tt.args.ctx, tt.args.limit, tt.args.nextCursor, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.GetWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WorkspaceService.GetWorkspaces() = %v, want %v", got, tt.want)
		}
	}

}

func TestWorkspaceService_GetWorkspaceWithChannels(t *testing.T) {
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
	workspaceService := NewWorkspaceService(dbService, redisService)
	workspace := Workspace{}
	workspace.ID = uint(2)
	workspace.UUID4 = []byte{28, 41, 191, 58, 70, 132, 73, 156, 165, 25, 44, 52, 138, 161, 50, 70}
	workspace.IDS = "1c29bf3a-4684-499c-a519-2c348aa13246"
	workspace.WorkspaceName = "Drive"
	workspace.WorkspaceDesc = "Drive"
	workspace.NumViews = uint(0)
	workspace.NumChannels = uint(1)
	workspace.Levelc = uint(1)
	workspace.ParentID = uint(1)
	workspace.NumChd = uint(0)
	workspace.UgroupID = uint(0)
	workspace.UserID = uint(1)
	workspace.Statusc = uint(1)
	workspace.CreatedAt = timeat
	workspace.UpdatedAt = timeat
	workspace.CreatedDay = uint(204)
	workspace.CreatedWeek = uint(30)
	workspace.CreatedMonth = uint(7)
	workspace.CreatedYear = uint(2019)
	workspace.UpdatedDay = uint(204)
	workspace.UpdatedWeek = uint(30)
	workspace.UpdatedMonth = uint(7)
	workspace.UpdatedYear = uint(2019)

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

	workspace.Channels = append(workspace.Channels, &channel)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		want    *Workspace
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				ID:        "1c29bf3a-4684-499c-a519-2c348aa13246",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &workspace,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetWorkspaceWithChannels(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.GetWorkspaceWithChannels() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WorkspaceService.GetWorkspaceWithChannels() = %v, want %v", got, tt.want)
		}
	}
}

func TestWorkspaceService_GetWorkspace(t *testing.T) {
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
	workspaceService := NewWorkspaceService(dbService, redisService)
	workspace := Workspace{}
	workspace.ID = uint(1)
	workspace.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	workspace.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	workspace.WorkspaceName = "Performance Portable Transmitter"
	workspace.WorkspaceDesc = "Performance Portable Transmitter"
	workspace.NumViews = uint(0)
	workspace.NumChannels = uint(0)
	workspace.Levelc = uint(0)
	workspace.ParentID = uint(0)
	workspace.NumChd = uint(1)
	workspace.UgroupID = uint(0)
	workspace.UserID = uint(1)
	workspace.Statusc = uint(1)
	workspace.CreatedAt = timeat
	workspace.UpdatedAt = timeat
	workspace.CreatedDay = uint(204)
	workspace.CreatedWeek = uint(30)
	workspace.CreatedMonth = uint(7)
	workspace.CreatedYear = uint(2019)
	workspace.UpdatedDay = uint(204)
	workspace.UpdatedWeek = uint(30)
	workspace.UpdatedMonth = uint(7)
	workspace.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		want    *Workspace
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				ID:        "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &workspace,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetWorkspace(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.GetWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WorkspaceService.GetWorkspace() = %v, want %v", got, tt.want)
		}
	}
}

func TestWorkspaceService_GetWorkspaceByID(t *testing.T) {
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
	workspaceService := NewWorkspaceService(dbService, redisService)
	workspace := Workspace{}
	workspace.ID = uint(1)
	workspace.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	workspace.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	workspace.WorkspaceName = "Performance Portable Transmitter"
	workspace.WorkspaceDesc = "Performance Portable Transmitter"
	workspace.NumViews = uint(0)
	workspace.NumChannels = uint(0)
	workspace.Levelc = uint(0)
	workspace.ParentID = uint(0)
	workspace.NumChd = uint(1)
	workspace.UgroupID = uint(0)
	workspace.UserID = uint(1)
	workspace.Statusc = uint(1)
	workspace.CreatedAt = timeat
	workspace.UpdatedAt = timeat
	workspace.CreatedDay = uint(204)
	workspace.CreatedWeek = uint(30)
	workspace.CreatedMonth = uint(7)
	workspace.CreatedYear = uint(2019)
	workspace.UpdatedDay = uint(204)
	workspace.UpdatedWeek = uint(30)
	workspace.UpdatedMonth = uint(7)
	workspace.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        uint
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		want    *Workspace
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				ID:        uint(1),
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &workspace,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetWorkspaceByID(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.GetWorkspaceByID() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WorkspaceService.GetWorkspaceByID() = %v, want %v", got, tt.want)
		}
	}
}

func TestWorkspaceService_GetTopLevelWorkspaces(t *testing.T) {
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
	workspaceService := NewWorkspaceService(dbService, redisService)
	workspace := Workspace{}
	workspaces := []*Workspace{}
	workspace.ID = uint(1)
	workspace.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	workspace.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	workspace.WorkspaceName = "Performance Portable Transmitter"
	workspace.WorkspaceDesc = "Performance Portable Transmitter"
	workspace.NumViews = uint(0)
	workspace.NumChannels = uint(0)
	workspace.Levelc = uint(0)
	workspace.ParentID = uint(0)
	workspace.NumChd = uint(1)
	workspace.UgroupID = uint(0)
	workspace.UserID = uint(1)
	workspace.Statusc = uint(1)
	workspace.CreatedAt = timeat
	workspace.UpdatedAt = timeat
	workspace.CreatedDay = uint(204)
	workspace.CreatedWeek = uint(30)
	workspace.CreatedMonth = uint(7)
	workspace.CreatedYear = uint(2019)
	workspace.UpdatedDay = uint(204)
	workspace.UpdatedWeek = uint(30)
	workspace.UpdatedMonth = uint(7)
	workspace.UpdatedYear = uint(2019)

	workspaces = append(workspaces, &workspace)
	type args struct {
		ctx       context.Context
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		want    []*Workspace
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    workspaces,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetTopLevelWorkspaces(tt.args.ctx, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.GetTopLevelWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WorkspaceService.GetTopLevelWorkspaces() = %v, want %v", got, tt.want)
		}
	}
}

func TestWorkspaceService_GetChildWorkspaces(t *testing.T) {
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
	workspaceService := NewWorkspaceService(dbService, redisService)
	workspace := Workspace{}
	workspaces := []*Workspace{}
	workspace.ID = uint(2)
	workspace.UUID4 = []byte{28, 41, 191, 58, 70, 132, 73, 156, 165, 25, 44, 52, 138, 161, 50, 70}
	workspace.IDS = "1c29bf3a-4684-499c-a519-2c348aa13246"
	workspace.WorkspaceName = "Drive"
	workspace.WorkspaceDesc = "Drive"
	workspace.NumViews = uint(0)
	workspace.NumChannels = uint(1)
	workspace.Levelc = uint(1)
	workspace.ParentID = uint(1)
	workspace.NumChd = uint(0)
	workspace.UgroupID = uint(0)
	workspace.UserID = uint(1)
	workspace.Statusc = uint(1)
	workspace.CreatedAt = timeat
	workspace.UpdatedAt = timeat
	workspace.CreatedDay = uint(204)
	workspace.CreatedWeek = uint(30)
	workspace.CreatedMonth = uint(7)
	workspace.CreatedYear = uint(2019)
	workspace.UpdatedDay = uint(204)
	workspace.UpdatedWeek = uint(30)
	workspace.UpdatedMonth = uint(7)
	workspace.UpdatedYear = uint(2019)

	workspaces = append(workspaces, &workspace)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		want    []*Workspace
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				ID:        "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    workspaces,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetChildWorkspaces(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.GetChildWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WorkspaceService.GetChildWorkspaces() = %v, want %v", got, tt.want)
		}
	}
}

func TestWorkspaceService_GetParentWorkspace(t *testing.T) {
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
	workspaceService := NewWorkspaceService(dbService, redisService)
	workspace := Workspace{}
	workspace.ID = uint(1)
	workspace.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	workspace.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	workspace.WorkspaceName = "Performance Portable Transmitter"
	workspace.WorkspaceDesc = "Performance Portable Transmitter"
	workspace.NumViews = uint(0)
	workspace.NumChannels = uint(0)
	workspace.Levelc = uint(0)
	workspace.ParentID = uint(0)
	workspace.NumChd = uint(1)
	workspace.UgroupID = uint(0)
	workspace.UserID = uint(1)
	workspace.Statusc = uint(1)
	workspace.CreatedAt = timeat
	workspace.UpdatedAt = timeat
	workspace.CreatedDay = uint(204)
	workspace.CreatedWeek = uint(30)
	workspace.CreatedMonth = uint(7)
	workspace.CreatedYear = uint(2019)
	workspace.UpdatedDay = uint(204)
	workspace.UpdatedWeek = uint(30)
	workspace.UpdatedMonth = uint(7)
	workspace.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		want    *Workspace
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				ID:        "1c29bf3a-4684-499c-a519-2c348aa13246",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &workspace,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetParentWorkspace(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.GetParentWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WorkspaceService.GetParentWorkspace() = %v, want %v", got, tt.want)
		}
	}
}

func TestWorkspaceService_UpdateWorkspace(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	workspaceService := NewWorkspaceService(dbService, redisService)
	form1 := Workspace{}
	form1.WorkspaceName = "Component Type B"
	form1.WorkspaceDesc = "Component Type B"

	type args struct {
		ctx       context.Context
		ID        string
		form      *Workspace
		UserID    string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				ID:        "1c29bf3a-4684-499c-a519-2c348aa13246",
				form:      &form1,
				UserID:    "29ea215b-8fb3-4453-b413-81a661e44495",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := tt.c.UpdateWorkspace(tt.args.ctx, tt.args.ID, tt.args.form, tt.args.UserID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.UpdateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestWorkspaceService_DeleteWorkspace(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	workspaceService := NewWorkspaceService(dbService, redisService)
	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *WorkspaceService
		args    args
		wantErr bool
	}{
		{
			c: workspaceService,
			args: args{
				ctx:       ctx,
				ID:        "1c29bf3a-4684-499c-a519-2c348aa13246",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := tt.c.DeleteWorkspace(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("WorkspaceService.DeleteWorkspace() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
