package main

import (
	"bot/pkg/memory"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	simplecron "github.com/sagleft/simple-cron"
)

type cronContainer struct {
	Cron      *simplecron.CronObject
	InProcess bool
}

func (c *cronContainer) markProcessing(isProcessing bool) {
	c.InProcess = isProcessing
}

type botCrons struct {
	ChannelContacts    *cronContainer
	FindChannels       *cronContainer
	RemoveOldRelations *cronContainer
	FindPeers          *cronContainer
	RemoveOldPeers     *cronContainer
}

func setupCronHandler(callback func(), timeout time.Duration, startImmediate bool) *cronContainer {
	c := &cronContainer{
		Cron: simplecron.NewCronHandler(callback, timeout),
	}
	go c.Cron.Run(startImmediate)
	return c
}

func (b *bot) checkStats() {
	channels, err := b.Memory.GetChannels()
	if err != nil {
		log.Println(err)
		return
	}

	for _, channel := range channels {
		b.Workers.CheckStats.AddEvent(checkChannelTask{
			Channel: channel,
		})
	}
}

func (b *bot) findChannels() {
	channels, err := b.Messenger.GetChannels()
	if err != nil {
		color.Red("get channels: %s", err.Error())
		return
	}

	for _, channel := range channels {
		if err := b.saveChannelIFNotExists(channel); err != nil {
			color.Red("save channel: %s", err.Error())
			return
		}
	}
}

func (b *bot) removeOldRelations() {
	relations, err := b.Memory.GetRelations()
	if err != nil {
		color.Red("get relations: %s", err.Error())
		return
	}

	for _, rel := range relations {
		if time.Since(rel.LastSeen) < maxRelationDuration {
			continue
		}

		if err := b.Memory.DeleteRelation(memory.ChannelUserRelation{
			ChannelID:      rel.ChannelID,
			UserPubkeyHash: rel.UserPubkeyHash,
		}); err != nil {
			color.Red("delete relation: %s", err.Error())
			return
		}
	}
}

func (b *bot) findPeers() {
	peers, err := b.Messenger.GetNetworkConnections()
	if err != nil {
		color.Red("get peers: %s", err.Error())
		return
	}

	for _, peer := range peers {
		if peer.Address == "" {
			color.Yellow("empty peer address given")
			return
		}

		peerAddress := strings.Split(peer.Address, ":")

		isExists, err := b.Memory.IsPeerExists(memory.Peer{
			IP: peerAddress[0],
		})
		if err != nil {
			color.Red("check peer exists in db: %s", err.Error())
			return
		}

		if !isExists {
			if err := b.Memory.SavePeer(memory.Peer{
				Direction: peer.Direction,
				IP:        peerAddress[0],
				Lat:       "", // TODO
				Lon:       "", // TODO
			}); err != nil {
				color.Red("save peer: %s", err.Error())
				return
			}
		}
	}
}

func (b *bot) removeOldPeers() {
	peers, err := b.Memory.GetPeers()
	if err != nil {
		color.Red("get peers: %s", err.Error())
		return
	}

	for _, peer := range peers {
		if time.Since(peer.CreatedAt) < maxPeerDuration {
			continue
		}

		if err := b.Memory.DeletePeer(memory.Peer{
			IP: peer.IP,
		}); err != nil {
			color.Red("delete peer: %s", err.Error())
			return
		}
	}
}
