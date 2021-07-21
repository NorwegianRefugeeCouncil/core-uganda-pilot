package pagination

import (
	"reflect"
	"testing"
)

func TestGetCurrentPage(t *testing.T) {
	type args struct {
		page int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should return first page if page number is missing",
			args: args{},
			want: 1,
		},
		{
			name: "should return first page if page number is zero",
			args: args{0},
			want: 1,
		},
		{
			name: "should return first page if page number is negative number",
			args: args{-1},
			want: 1,
		},
		{
			name: "should return first page",
			args: args{1},
			want: 1,
		},
		{
			name: "should return page numbers bigger than one",
			args: args{2},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCurrentPage(tt.args.page); got != tt.want {
				t.Errorf("GetCurrentPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMaxPerPage(t *testing.T) {
	type args struct {
		perPage int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should return 100 records by default if perPage is missing",
			args: args{},
			want: 100,
		},
		{
			name: "should return 100 records by default if perPage is zero",
			args: args{0},
			want: 100,
		},
		{
			name: "should return 100 records by default if perPage is negative number",
			args: args{-1},
			want: 100,
		},
		{
			name: "should return max 250 records if perPage is bigger than 250",
			args: args{251},
			want: 250,
		},
		{
			name: "should return perPage value if number between zero and 250",
			args: args{42},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaxPerPage(tt.args.perPage); got != tt.want {
				t.Errorf("GetMaxPerPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLastPageCount(t *testing.T) {
	type args struct {
		totalCount int
		maxPerPage int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should return 1 if totalCount is less than maxPerPage",
			args: args{1, 2},
			want: 1,
		},
		{
			name: "should return 1 if totalCount is eaqual to maxPerPage",
			args: args{2, 2},
			want: 1,
		},
		{
			name: "should divide totalCount by maxPerPage",
			args: args{10, 2},
			want: 5,
		},
		{
			name: "should divide totalCount by maxPerPage and return ceil value",
			args: args{11, 2},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLastPageCount(tt.args.totalCount, tt.args.maxPerPage); got != tt.want {
				t.Errorf("getLastPageCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSortOptionType(t *testing.T) {
	type args struct {
		sort string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should default to ascending if sort type is missing",
			args: args{},
			want: 1,
		},
		{
			name: "should default to ascending if sort type is empty or unknown",
			args: args{""},
			want: 1,
		},
		{
			name: "should return ascending if sort type is asc",
			args: args{"asc"},
			want: 1,
		},
		{
			name: "should return -1 if sort is descending",
			args: args{"desc"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSortOptionType(tt.args.sort); got != tt.want {
				t.Errorf("GetSortOptionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSortDirection(t *testing.T) {
	type args struct {
		sort string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should default to ascending if sort string is missing",
			args: args{},
			want: "asc",
		},
		{
			name: "should default to ascending if sort string is empty or unknown",
			args: args{""},
			want: "asc",
		},
		{
			name: "should return asc if sort string is ascending",
			args: args{"asc"},
			want: "asc",
		},
		{
			name: "should return desc if sort string is descending",
			args: args{"desc"},
			want: "desc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSortDirection(tt.args.sort); got != tt.want {
				t.Errorf("getSortDirection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLinks(t *testing.T) {
	type args struct {
		searchParam  string
		currentPage  int
		pageCount    int
		perPageParam string
		sortParam    string
	}
	tests := []struct {
		name string
		args args
		want Links
	}{
		{
			name: "should generate params for pagination link buttons",
			args: args{
				searchParam:  "searchParam=abc",
				currentPage:  3,
				pageCount:    5,
				perPageParam: "&perPage=1",
				sortParam:    "&sort=asc",
			},
			want: Links{
				First:    "?searchParam=abc&page=1&perPage=1&sort=asc",
				Previous: "?searchParam=abc&page=2&perPage=1&sort=asc",
				Self:     "?searchParam=abc&page=3&perPage=1&sort=asc",
				Next:     "?searchParam=abc&page=4&perPage=1&sort=asc",
				Last:     "?searchParam=abc&page=5&perPage=1&sort=asc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLinks(tt.args.searchParam, tt.args.currentPage, tt.args.pageCount, tt.args.perPageParam, tt.args.sortParam); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPaginationMetaData(t *testing.T) {
	type args struct {
		totalCount  int
		currentPage int
		perPage     int
		sortValue   string
		searchValue string
	}
	tests := []struct {
		name string
		args args
		want Pagination
	}{
		{
			name: "should generate pagination meta data",
			args: args{
				totalCount:  6,
				currentPage: 1,
				perPage:     2,
				sortValue:   "desc",
				searchValue: "abc",
			},
			want: Pagination{
				Page:       1,
				PerPage:    2,
				PageCount:  3,
				TotalCount: 6,
				Sort:       "desc",
				Links: Links{
					First:    "?searchParam=abc&page=1&perPage=2&sort=desc",
					Previous: "?searchParam=abc&page=1&perPage=2&sort=desc",
					Self:     "?searchParam=abc&page=1&perPage=2&sort=desc",
					Next:     "?searchParam=abc&page=2&perPage=2&sort=desc",
					Last:     "?searchParam=abc&page=3&perPage=2&sort=desc",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPaginationMetaData(tt.args.totalCount, tt.args.currentPage, tt.args.perPage, tt.args.sortValue, tt.args.searchValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPaginationMetaData() = %v, want %v", got, tt.want)
			}
		})
	}
}
