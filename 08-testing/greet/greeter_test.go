package utils

import (
	"testing"
	"time"
)

func Test_Greeter_4_Morning(t *testing.T) {
	// Arrange
	userName := "Magesh"
	expected := "Hi Magesh, Good Morning!"

	// creating an instance of the Mock type created by Mockery
	timeServiceMock := NewMockTimeProvider(t)

	// Configure the Mock to return our hardcoded time when "GetCurrent" method called
	timeServiceMock.On("GetCurrent").Return(time.Date(2025, time.April, 25, 10, 0, 0, 0, time.UTC))

	sut := NewGreeter(userName, timeServiceMock)
	// Act
	actual := sut.Greet()

	// Assert
	if expected != actual {
		t.Errorf("expected : %q, actual : %q\n", expected, actual)
	}
}

func Test_Greeter_4_After_Morning(t *testing.T) {
	// Arrange
	userName := "Magesh"
	expected := "Hi Magesh, Good Day!"
	timeServiceMock := NewMockTimeProvider(t)

	timeServiceMock.On("GetCurrent").Return(time.Date(2025, time.April, 25, 15, 0, 0, 0, time.UTC))

	sut := NewGreeter(userName, timeServiceMock)
	// Act
	actual := sut.Greet()

	// Assert
	if expected != actual {
		t.Errorf("expected : %q, actual : %q\n", expected, actual)
	}
}
