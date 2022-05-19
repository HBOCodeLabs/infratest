package cassandra

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hbocodelabs/infratest/mock"
	"github.com/stretchr/testify/assert"
)

func TestAssertCassandraQuerySucceedsE_Success(t *testing.T) {

	ctrl := gomock.NewController(t)
	query := mock.NewMockQueryInterface(ctrl)

	ctx := context.Background()
	fakeTest := &testing.T{}

	query.EXPECT().Exec().Times(1).DoAndReturn(func() (err error) {
		return nil
	})
	AssertCassandraQuerySucceeds(fakeTest, ctx, query)
	assert.False(t, fakeTest.Failed())
}

func TestAssertCassandraQuerySucceedsE_Failed(t *testing.T) {

	ctrl := gomock.NewController(t)
	query := mock.NewMockQueryInterface(ctrl)

	ctx := context.Background()
	fakeTest := &testing.T{}

	query.EXPECT().Exec().Times(1).DoAndReturn(func() (err error) {
		return fmt.Errorf("I am an Error")
	})
	AssertCassandraQuerySucceeds(fakeTest, ctx, query)
	assert.True(t, fakeTest.Failed())
}
