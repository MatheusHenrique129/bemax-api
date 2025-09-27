package apierrors_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
	"github.com/stretchr/testify/assert"
)

func TestRestError(t *testing.T) {
	t.Run("NewRestError", func(t *testing.T) {
		stringError := "error generic test"
		var causeList apierrors.CauseList
		apiError := http.StatusBadRequest

		err := apierrors.NewRestError("test", stringError, apiError, causeList)

		assert.Equal(t, err.Status(), apiError)
		assert.Equal(t, err.Code(), stringError)
		assert.Equal(t, err.Cause(), causeList)
		assert.NotNil(t, err.Message())
	})

	t.Run("NewNotFoundRestError", func(t *testing.T) {
		err := apierrors.NewNotFoundRestError("test")

		assert.Equal(t, err.Status(), http.StatusNotFound)
		assert.Nil(t, err.Cause())
		assert.NotNil(t, err.Code())
		assert.NotNil(t, err.Message())
	})

	t.Run("NewBadRequestRestError", func(t *testing.T) {
		err := apierrors.NewBadRequestRestError("test")

		assert.Equal(t, err.Status(), http.StatusBadRequest)
		assert.Nil(t, err.Cause())
		assert.NotNil(t, err.Code())
		assert.NotNil(t, err.Message())
	})

	t.Run("NewBadRequestValidationRestError", func(t *testing.T) {
		stringError := "error generic test"
		var causeList apierrors.CauseList

		err := apierrors.NewBadRequestValidationRestError("test", stringError, causeList)

		assert.Equal(t, err.Status(), http.StatusBadRequest)
		assert.Equal(t, err.Code(), stringError)
		assert.Equal(t, err.Cause(), causeList)
		assert.NotNil(t, err.Message())
	})

	t.Run("NewUnauthorizedRestError", func(t *testing.T) {
		err := apierrors.NewUnauthorizedRestError("test")

		assert.Equal(t, err.Status(), http.StatusUnauthorized)
		assert.Nil(t, err.Cause())
		assert.NotNil(t, err.Code())
		assert.NotNil(t, err.Message())
	})

	t.Run("NewInternalServerRestError", func(t *testing.T) {
		err := apierrors.NewInternalServerRestError("test", errors.New("error generic test"))

		assert.Equal(t, err.Status(), http.StatusInternalServerError)
		assert.NotNil(t, err.Code())
		assert.NotNil(t, err.Cause())
		assert.NotNil(t, err.Message())
		assert.NotNil(t, err.Error())
		assert.NotNil(t, err.Cause())
	})
}
