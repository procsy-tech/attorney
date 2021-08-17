package service

import (
	"testing"
    "github.com/procsy-tech/attorney/dto"
	"github.com/procsy-tech/attorney/repository"
	"github.com/kbkontrakt/hlfabric-ccdevkit/logs"
	. "github.com/smartystreets/goconvey/convey"

	gomock "github.com/golang/mock/gomock"
)

func TestPOAServiceCreate(t *testing.T) {
	Convey("POA Create", t, func(c C) {
		// prepare dummy service .

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewPOAServiceImpl(
			logs.DummyLogger(),
			repository.NewMockRepository(ctrl),
		)

		c.Convey("Given POAService", func(c C) {
			c.Convey("When invoking method Create", func(c C) {
				// test logic needs to be implemented .
				
				var (
					request    = &dto.CreateRequest{}
				)
    			c.Convey("It should return error", func(c C) {
					_, err := svc.Create(request.POA)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
func TestPOAServiceConfirmAttorney(t *testing.T) {
	Convey("POA ConfirmAttorney", t, func(c C) {
		// prepare dummy service .

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewPOAServiceImpl(
			logs.DummyLogger(),
			repository.NewMockRepository(ctrl),
		)

		c.Convey("Given POAService", func(c C) {
			c.Convey("When invoking method ConfirmAttorney", func(c C) {
				// test logic needs to be implemented .
				
				var (
					request    = &dto.ConfirmAttorneyRequest{}
				)
    			c.Convey("It should return error", func(c C) {
					err := svc.ConfirmAttorney(request.ID)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

