#!/bin/sh

fsrpl dump "calendars/*" --path dump/calendars
fsrpl dump "expoPushTokens/*" --path dump/expoPushTokens
fsrpl dump "items/*" --path dump/items
fsrpl dump "itemDetails/*" --path dump/itemDetails
fsrpl dump "plans/*" --path dump/plans
fsrpl dump "preResisterItems/*" --path dump/preResisterItems
fsrpl dump "userIntegrations/*" --path dump/userIntegrations
fsrpl dump "users/*" --path dump/users

