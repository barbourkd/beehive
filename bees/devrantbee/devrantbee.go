/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Nicolas Martin <penguwingithub@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package devrantbee is a Bee that can post blogs & quotes on Devrant.
package devrantbee

import (
	"github.com/jayeshsolanki93/devgorant"

	"github.com/muesli/beehive/bees"
)

// DevrantBee is a Bee that can post blogs & quotes on Devrant.
type DevrantBee struct {
	bees.Bee

	client *devgorant.Client

	evchan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *DevrantBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "surprise":
		rant, err := mod.client.Surprise()
		if err != nil {
			mod.LogErrorf("Failed to fetch surprise rant: %v", err)
		}

		mod.triggerEvent(rant)

	case "weekly":
		rants, err := mod.client.WeeklyRants()
		if err != nil {
			mod.LogErrorf("Failed to fetch weekly rants: %v", err)
		}

		for _, v := range rants {
			mod.triggerEvent(v)
		}

	case "rant":
		limit := 0
		action.Options.Bind("limit", &limit)

		rants, err := mod.client.Rants("", limit, 0)
		if err != nil {
			mod.LogErrorf("Failed to fetch rants: %v", err)
		}

		for i := range rants {
			mod.triggerEvent(rants[i])
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *DevrantBee) triggerEvent(rant devgorant.RantModel) {
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "rant",
		Options: []bees.Placeholder{
			{
				Name:  "ID",
				Type:  "int",
				Value: rant.Id,
			},
			{
				Name:  "text",
				Type:  "string",
				Value: rant.Text,
			},
			{
				Name:  "upvotes",
				Type:  "int",
				Value: rant.Upvotes,
			},
			{
				Name:  "downvotes",
				Type:  "int",
				Value: rant.Downvotes,
			},
			{
				Name:  "score",
				Type:  "int",
				Value: rant.Score,
			},
			{
				Name:  "created_time",
				Type:  "int",
				Value: rant.CreatedTime,
			},
			{
				Name:  "num_comments",
				Type:  "int",
				Value: rant.NumComments,
			},
			{
				Name:  "user_id",
				Type:  "int",
				Value: rant.UserId,
			},
			{
				Name:  "username",
				Type:  "string",
				Value: rant.UserUsername,
			},
			{
				Name:  "user_score",
				Type:  "int",
				Value: rant.UserScore,
			},
		},
	}
	mod.evchan <- ev
}

// Run executes the Bee's event loop.
func (mod *DevrantBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	// Setting up the client, unfortunately we can't really log in as an user
	mod.client = devgorant.New()
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *DevrantBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}
