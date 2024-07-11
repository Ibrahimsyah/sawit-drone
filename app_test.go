package main

import (
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUtil := NewMockUtilProvider(ctrl)

	got := NewApp(mockUtil)
	assert.Equal(t, &App{util: mockUtil}, got)
}

func TestApp_Start(t *testing.T) {
	type mockFields struct {
		util *MockUtilProvider
	}

	tests := []struct {
		name       string
		mock       func(mock mockFields)
		wantOsExit int
	}{
		{
			name:       "given invalid input then expect os.Exit 1",
			wantOsExit: 1,
			mock: func(mock mockFields) {
				length, width, count := 0, 0, 0
				mock.util.EXPECT().Scanln(&length, &width, &count).Do(
					func(target ...any) {
						*target[0].(*int) = 1
						*target[1].(*int) = 2
						*target[2].(*int) = 50001
					},
				)
				osExit = func(code int) {
					panic(code)
				}
			},
		},
		{
			name:       "given invalid tree height then expect os.Exit 1",
			wantOsExit: 1,
			mock: func(mock mockFields) {
				length, width, count := 0, 0, 0
				mock.util.EXPECT().Scanln(&length, &width, &count).Do(
					func(target ...any) {
						*target[0].(*int) = 1
						*target[1].(*int) = 2
						*target[2].(*int) = 1
					},
				)

				x, y, height := 0, 0, 0
				mock.util.EXPECT().Scanln(&x, &y, &height).Do(
					func(target ...any) {
						*target[0].(*int) = 1
						*target[1].(*int) = 2
						*target[2].(*int) = 31
					},
				)
				osExit = func(code int) {
					panic(code)
				}
			},
		},
		{
			name: "given success then expect no error and panic",
			mock: func(mock mockFields) {
				length, width, count := 0, 0, 0
				mock.util.EXPECT().Scanln(&length, &width, &count).Do(
					func(target ...any) {
						*target[0].(*int) = 1
						*target[1].(*int) = 2
						*target[2].(*int) = 1
					},
				)
				x, y, height := 0, 0, 0
				mock.util.EXPECT().Scanln(&x, &y, &height).Do(
					func(target ...any) {
						*target[0].(*int) = 5
						*target[1].(*int) = 2
						*target[2].(*int) = 2
					},
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Rollback the mock
			defer func() {
				osExit = os.Exit
			}()

			ctrl := gomock.NewController(t)
			mockFields := mockFields{
				util: NewMockUtilProvider(ctrl),
			}
			test.mock(mockFields)

			app := &App{
				util: mockFields.util,
			}
			if test.wantOsExit != 0 {
				assert.PanicsWithValue(t, test.wantOsExit, app.Start)
			} else {
				app.Start()
			}
		})
	}
}

func TestApp_calculateFlyDistance(t *testing.T) {
	type args struct {
		length  int
		width   int
		treeMap map[string]int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "given no tree then return the horizontal distance + landing + take off",
			args: args{
				length:  10,
				width:   1,
				treeMap: make(map[string]int),
			},
			want: 92,
		},
		{
			name: "given 5 x 1 field with 3 tree on the middle then return 54",
			args: args{
				length: 5,
				width:  1,
				treeMap: map[string]int{
					"2,1": 5,
					"3,1": 3,
					"4,1": 4,
				},
			},
			want: 54,
		},
		{
			name: "given 5 x 1 field with 3 unordered tree on the middle then return 54",
			args: args{
				length: 5,
				width:  1,
				treeMap: map[string]int{
					"3,1": 3,
					"4,1": 4,
					"2,1": 5,
				},
			},
			want: 54,
		},
		{
			name: "given 5 x 1 field with 3 tree with same height on the middle then return 62",
			args: args{
				length: 5,
				width:  1,
				treeMap: map[string]int{
					"3,1": 10,
					"4,1": 10,
					"2,1": 10,
				},
			},
			want: 62,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := &App{}
			got := app.calculateFlyDistance(test.args.length, test.args.width, test.args.treeMap)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_calculateHorizontalDistance(t *testing.T) {
	type args struct {
		length int
		width  int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "given 5x1 field then return 4",
			args: args{
				length: 5,
				width:  1,
			},
			want: 4,
		},
		{
			name: "given 5x2 field then return 9",
			args: args{
				length: 5,
				width:  2,
			},
			want: 9,
		},
		{
			name: "given 5x3 field then return 14",
			args: args{
				length: 5,
				width:  3,
			},
			want: 14,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := &App{}
			got := app.calculateHorizontalDistance(test.args.length, test.args.width)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_getNextPlotCoordinate(t *testing.T) {
	type args struct {
		length int
		x      int
		y      int
	}
	tests := []struct {
		name  string
		args  args
		wantX int
		wantY int
	}{
		{
			name: "given x = 1 then return x + 1, y",
			args: args{
				length: 4,
				x:      1,
				y:      1,
			},
			wantX: 2,
			wantY: 1,
		},
		{
			name: "given x = 1 and y = 2 then return x, y + 1",
			args: args{
				length: 4,
				x:      1,
				y:      2,
			},
			wantX: 1,
			wantY: 3,
		},
		{
			name: "given x = length then return x, y + 1",
			args: args{
				length: 4,
				x:      4,
				y:      1,
			},
			wantX: 4,
			wantY: 2,
		},
		{
			name: "given x = length and y = 2 then return x - 1, y",
			args: args{
				length: 4,
				x:      4,
				y:      2,
			},
			wantX: 3,
			wantY: 2,
		},
		{
			name: "given x not 1 and not equal to length and y even then return x - 1, y",
			args: args{
				length: 4,
				x:      2,
				y:      2,
			},
			wantX: 1,
			wantY: 2,
		},
		{
			name: "given x not 1 and not equal to length and y odd then return x + 1, y",
			args: args{
				length: 4,
				x:      2,
				y:      1,
			},
			wantX: 3,
			wantY: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := &App{}
			gotX, gotY := app.getNextPlotCoordinate(test.args.length, test.args.x, test.args.y)
			assert.Equal(t, test.wantX, gotX)
			assert.Equal(t, test.wantY, gotY)
		})
	}
}
