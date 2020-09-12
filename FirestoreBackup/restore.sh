#!/bin/sh

fsrpl restore "calendars/*" --path dump/calendars
fsrpl restore "expoPushTokens/*" --path dump/expoPushTokens
fsrpl restore "items/*" --path dump/items
fsrpl restore "itemDetails/*" --path dump/itemDetails
fsrpl restore "plans/*" --path dump/plans
fsrpl restore "preResisterItems/*" --path dump/preResisterItems
fsrpl restore "userIntegrations/*" --path dump/userIntegrations
fsrpl restore "users/*" --path dump/users

