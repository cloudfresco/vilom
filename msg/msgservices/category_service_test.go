package msgservices

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/cloudfresco/vilom/testhelpers"
	_ "github.com/go-sql-driver/mysql"
)

func TestCategoryService_GetCategories(t *testing.T) {
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
	categoryService := NewCategoryService(dbService, redisService)
	cat := Category{}
	cats := []*Category{}
	cat.ID = uint(1)
	cat.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	cat.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	cat.CategoryName = "Performance Portable Transmitter"
	cat.CategoryDesc = "Performance Portable Transmitter"
	cat.NumViews = uint(0)
	cat.NumTopics = uint(0)
	cat.Levelc = uint(0)
	cat.ParentID = uint(0)
	cat.NumChd = uint(1)
	cat.UgroupID = uint(0)
	cat.UserID = uint(1)
	cat.Statusc = uint(1)
	cat.CreatedAt = timeat
	cat.UpdatedAt = timeat
	cat.CreatedDay = uint(204)
	cat.CreatedWeek = uint(30)
	cat.CreatedMonth = uint(7)
	cat.CreatedYear = uint(2019)
	cat.UpdatedDay = uint(204)
	cat.UpdatedWeek = uint(30)
	cat.UpdatedMonth = uint(7)
	cat.UpdatedYear = uint(2019)

	cats = append(cats, &cat)
	nextc := "MA=="
	c1 := CategoryCursor{cats, nextc}

	type args struct {
		ctx        context.Context
		limit      string
		nextCursor string
		userEmail  string
		requestID  string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		want    *CategoryCursor
		wantErr bool
	}{
		{
			c: categoryService,
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
		got, err := tt.c.GetCategories(tt.args.ctx, tt.args.limit, tt.args.nextCursor, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.GetCategories() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CategoryService.GetCategories() = %v, want %v", got, tt.want)
		}
	}

}

func TestCategoryService_GetCategoryWithTopics(t *testing.T) {
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
	categoryService := NewCategoryService(dbService, redisService)
	cat := Category{}
	cat.ID = uint(2)
	cat.UUID4 = []byte{28, 41, 191, 58, 70, 132, 73, 156, 165, 25, 44, 52, 138, 161, 50, 70}
	cat.IDS = "1c29bf3a-4684-499c-a519-2c348aa13246"
	cat.CategoryName = "Drive"
	cat.CategoryDesc = "Drive"
	cat.NumViews = uint(0)
	cat.NumTopics = uint(1)
	cat.Levelc = uint(1)
	cat.ParentID = uint(1)
	cat.NumChd = uint(0)
	cat.UgroupID = uint(0)
	cat.UserID = uint(1)
	cat.Statusc = uint(1)
	cat.CreatedAt = timeat
	cat.UpdatedAt = timeat
	cat.CreatedDay = uint(204)
	cat.CreatedWeek = uint(30)
	cat.CreatedMonth = uint(7)
	cat.CreatedYear = uint(2019)
	cat.UpdatedDay = uint(204)
	cat.UpdatedWeek = uint(30)
	cat.UpdatedMonth = uint(7)
	cat.UpdatedYear = uint(2019)

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

	cat.Topics = append(cat.Topics, &topic)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		want    *Category
		wantErr bool
	}{
		{
			c: categoryService,
			args: args{
				ctx:       ctx,
				ID:        "1c29bf3a-4684-499c-a519-2c348aa13246",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &cat,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetCategoryWithTopics(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.GetCategoryWithTopics() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CategoryService.GetCategoryWithTopics() = %v, want %v", got, tt.want)
		}
	}
}

func TestCategoryService_GetCategory(t *testing.T) {
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
	categoryService := NewCategoryService(dbService, redisService)
	cat := Category{}
	cat.ID = uint(1)
	cat.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	cat.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	cat.CategoryName = "Performance Portable Transmitter"
	cat.CategoryDesc = "Performance Portable Transmitter"
	cat.NumViews = uint(0)
	cat.NumTopics = uint(0)
	cat.Levelc = uint(0)
	cat.ParentID = uint(0)
	cat.NumChd = uint(1)
	cat.UgroupID = uint(0)
	cat.UserID = uint(1)
	cat.Statusc = uint(1)
	cat.CreatedAt = timeat
	cat.UpdatedAt = timeat
	cat.CreatedDay = uint(204)
	cat.CreatedWeek = uint(30)
	cat.CreatedMonth = uint(7)
	cat.CreatedYear = uint(2019)
	cat.UpdatedDay = uint(204)
	cat.UpdatedWeek = uint(30)
	cat.UpdatedMonth = uint(7)
	cat.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		want    *Category
		wantErr bool
	}{
		{
			c: categoryService,
			args: args{
				ctx:       ctx,
				ID:        "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &cat,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetCategory(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.GetCategory() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CategoryService.GetCategory() = %v, want %v", got, tt.want)
		}
	}
}

func TestCategoryService_GetCategoryByID(t *testing.T) {
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
	categoryService := NewCategoryService(dbService, redisService)
	cat := Category{}
	cat.ID = uint(1)
	cat.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	cat.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	cat.CategoryName = "Performance Portable Transmitter"
	cat.CategoryDesc = "Performance Portable Transmitter"
	cat.NumViews = uint(0)
	cat.NumTopics = uint(0)
	cat.Levelc = uint(0)
	cat.ParentID = uint(0)
	cat.NumChd = uint(1)
	cat.UgroupID = uint(0)
	cat.UserID = uint(1)
	cat.Statusc = uint(1)
	cat.CreatedAt = timeat
	cat.UpdatedAt = timeat
	cat.CreatedDay = uint(204)
	cat.CreatedWeek = uint(30)
	cat.CreatedMonth = uint(7)
	cat.CreatedYear = uint(2019)
	cat.UpdatedDay = uint(204)
	cat.UpdatedWeek = uint(30)
	cat.UpdatedMonth = uint(7)
	cat.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        uint
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		want    *Category
		wantErr bool
	}{
		{
			c: categoryService,
			args: args{
				ctx:       ctx,
				ID:        uint(1),
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &cat,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetCategoryByID(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.GetCategoryByID() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CategoryService.GetCategoryByID() = %v, want %v", got, tt.want)
		}
	}
}

func TestCategoryService_GetTopLevelCategories(t *testing.T) {
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
	categoryService := NewCategoryService(dbService, redisService)
	cat := Category{}
	cats := []*Category{}
	cat.ID = uint(1)
	cat.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	cat.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	cat.CategoryName = "Performance Portable Transmitter"
	cat.CategoryDesc = "Performance Portable Transmitter"
	cat.NumViews = uint(0)
	cat.NumTopics = uint(0)
	cat.Levelc = uint(0)
	cat.ParentID = uint(0)
	cat.NumChd = uint(1)
	cat.UgroupID = uint(0)
	cat.UserID = uint(1)
	cat.Statusc = uint(1)
	cat.CreatedAt = timeat
	cat.UpdatedAt = timeat
	cat.CreatedDay = uint(204)
	cat.CreatedWeek = uint(30)
	cat.CreatedMonth = uint(7)
	cat.CreatedYear = uint(2019)
	cat.UpdatedDay = uint(204)
	cat.UpdatedWeek = uint(30)
	cat.UpdatedMonth = uint(7)
	cat.UpdatedYear = uint(2019)

	cats = append(cats, &cat)
	type args struct {
		ctx       context.Context
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		want    []*Category
		wantErr bool
	}{
		{
			c: categoryService,
			args: args{
				ctx:       ctx,
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    cats,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetTopLevelCategories(tt.args.ctx, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.GetTopLevelCategories() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CategoryService.GetTopLevelCategories() = %v, want %v", got, tt.want)
		}
	}
}

func TestCategoryService_GetChildCategories(t *testing.T) {
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
	categoryService := NewCategoryService(dbService, redisService)
	cat := Category{}
	cats := []*Category{}
	cat.ID = uint(2)
	cat.UUID4 = []byte{28, 41, 191, 58, 70, 132, 73, 156, 165, 25, 44, 52, 138, 161, 50, 70}
	cat.IDS = "1c29bf3a-4684-499c-a519-2c348aa13246"
	cat.CategoryName = "Drive"
	cat.CategoryDesc = "Drive"
	cat.NumViews = uint(0)
	cat.NumTopics = uint(1)
	cat.Levelc = uint(1)
	cat.ParentID = uint(1)
	cat.NumChd = uint(0)
	cat.UgroupID = uint(0)
	cat.UserID = uint(1)
	cat.Statusc = uint(1)
	cat.CreatedAt = timeat
	cat.UpdatedAt = timeat
	cat.CreatedDay = uint(204)
	cat.CreatedWeek = uint(30)
	cat.CreatedMonth = uint(7)
	cat.CreatedYear = uint(2019)
	cat.UpdatedDay = uint(204)
	cat.UpdatedWeek = uint(30)
	cat.UpdatedMonth = uint(7)
	cat.UpdatedYear = uint(2019)

	cats = append(cats, &cat)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		want    []*Category
		wantErr bool
	}{
		{
			c: categoryService,
			args: args{
				ctx:       ctx,
				ID:        "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    cats,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetChildCategories(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.GetChildCategories() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CategoryService.GetChildCategories() = %v, want %v", got, tt.want)
		}
	}
}

func TestCategoryService_GetParentCategory(t *testing.T) {
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
	categoryService := NewCategoryService(dbService, redisService)
	cat := Category{}
	cat.ID = uint(1)
	cat.UUID4 = []byte{27, 209, 136, 138, 219, 254, 69, 16, 167, 173, 169, 143, 105, 253, 10, 107}
	cat.IDS = "1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b"
	cat.CategoryName = "Performance Portable Transmitter"
	cat.CategoryDesc = "Performance Portable Transmitter"
	cat.NumViews = uint(0)
	cat.NumTopics = uint(0)
	cat.Levelc = uint(0)
	cat.ParentID = uint(0)
	cat.NumChd = uint(1)
	cat.UgroupID = uint(0)
	cat.UserID = uint(1)
	cat.Statusc = uint(1)
	cat.CreatedAt = timeat
	cat.UpdatedAt = timeat
	cat.CreatedDay = uint(204)
	cat.CreatedWeek = uint(30)
	cat.CreatedMonth = uint(7)
	cat.CreatedYear = uint(2019)
	cat.UpdatedDay = uint(204)
	cat.UpdatedWeek = uint(30)
	cat.UpdatedMonth = uint(7)
	cat.UpdatedYear = uint(2019)

	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		want    *Category
		wantErr bool
	}{
		{
			c: categoryService,
			args: args{
				ctx:       ctx,
				ID:        "1c29bf3a-4684-499c-a519-2c348aa13246",
				userEmail: "abcd145@gmail.com",
				requestID: "bks1m1g91jau4nkks2f0",
			},
			want:    &cat,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.c.GetParentCategory(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.GetParentCategory() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CategoryService.GetParentCategory() = %v, want %v", got, tt.want)
		}
	}
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	categoryService := NewCategoryService(dbService, redisService)
	form1 := Category{}
	form1.CategoryName = "Component Type B"
	form1.CategoryDesc = "Component Type B"

	type args struct {
		ctx       context.Context
		ID        string
		form      *Category
		UserID    string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		wantErr bool
	}{
		{
			c: categoryService,
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
		if err := tt.c.UpdateCategory(tt.args.ctx, tt.args.ID, tt.args.form, tt.args.UserID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.UpdateCategory() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	categoryService := NewCategoryService(dbService, redisService)
	type args struct {
		ctx       context.Context
		ID        string
		userEmail string
		requestID string
	}
	tests := []struct {
		c       *CategoryService
		args    args
		wantErr bool
	}{
		{
			c: categoryService,
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
		if err := tt.c.DeleteCategory(tt.args.ctx, tt.args.ID, tt.args.userEmail, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("CategoryService.DeleteCategory() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
