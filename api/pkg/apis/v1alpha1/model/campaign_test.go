/*

	MIT License

	Copyright (c) Microsoft Corporation.

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE

*/

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCampaignMatch(t *testing.T) {
	campaign1 := CampaignSpec{
		Name: "name",
	}
	campaign2 := CampaignSpec{
		Name: "name",
	}
	equal, err := campaign1.DeepEquals(campaign2)
	assert.Nil(t, err)
	assert.True(t, equal)
}

func TestCampaignMatchOneEmpty(t *testing.T) {
	campaign1 := CampaignSpec{
		Name: "name",
	}
	res, err := campaign1.DeepEquals(nil)
	assert.Errorf(t, err, "parameter is not a CampaignSpec type")
	assert.False(t, res)
}

func TestCampaignRoleNotMatch(t *testing.T) {
	campaign1 := CampaignSpec{
		Name: "name",
	}
	campaign2 := CampaignSpec{
		Name: "name1",
	}
	equal, err := campaign1.DeepEquals(campaign2)
	assert.Nil(t, err)
	assert.False(t, equal)
}
func TestScheduleTimeZone(t *testing.T) {
	schedule := ScheduleSpec{
		Date: "2020-01-01",
		Time: "12:00PM",
		Zone: "PST",
	}
	dt, err := schedule.GetTime()
	assert.Nil(t, err)
	assert.Equal(t, "2020-01-01 12:00:00 -0800 PST", dt.String())
}

func TestScheduleIANATimeZone(t *testing.T) {
	schedule := ScheduleSpec{
		Date: "2020-01-01",
		Time: "12:00PM",
		Zone: "America/Los_Angeles",
	}
	dt, err := schedule.GetTime()
	assert.Nil(t, err)
	assert.Equal(t, "2020-01-01 12:00:00 -0800 PST", dt.String())
}

func TestScheduleEmpty(t *testing.T) {
	schedule := ScheduleSpec{
		Date: "2020-01-01",
		Time: "12:00PM",
		Zone: "", //this is equivalent to UTC
	}
	dt, err := schedule.GetTime()
	assert.Nil(t, err)
	assert.Equal(t, "2020-01-01 12:00:00 +0000 UTC", dt.String())
}

func TestScheduleUTC(t *testing.T) {
	schedule := ScheduleSpec{
		Date: "2020-01-01",
		Time: "12:00PM",
		Zone: "UTC",
	}
	dt, err := schedule.GetTime()
	assert.Nil(t, err)
	assert.Equal(t, "2020-01-01 12:00:00 +0000 UTC", dt.String())
}

// TODO: This test works only in PST timezone, need to fix it for all time zones
// func TestScheduleLocal(t *testing.T) {
// 	schedule := ScheduleSpec{
// 		Date: "2020-01-01",
// 		Time: "12:00PM",
// 		Zone: "Local",
// 	}
// 	dt, err := schedule.GetTime()
// 	assert.Nil(t, err)
// 	assert.Equal(t, "2020-01-01 12:00:00 -0800 PST", dt.String())
// }
