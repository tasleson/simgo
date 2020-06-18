// SPDX-License-Identifier: 0BSD

package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	lsm "github.com/libstorage/libstoragemgmt-golang"
	"github.com/libstorage/libstoragemgmt-golang/errors"
)

type pluginData struct {
	reg lsm.PluginRegister
	c   *lsm.ClientConnection
	tmo uint32
}

var state pluginData

func register(p *lsm.PluginRegister) error {

	parsed, pE := url.Parse(p.URI)
	if pE != nil {
		return pE
	}

	values := strings.Split(parsed.RawQuery, "=")
	if len(values) != 2 || values[0] != "forward" {
		return &errors.LsmError{
			Code:    errors.InvalidArgument,
			Message: fmt.Sprintf("expected query string to be 'forward=<otherplugin>' got %s", parsed.RawQuery)}
	}

	var cE error
	state.c, cE = lsm.Client(fmt.Sprintf("%s://", values[1]), "", p.Timeout)
	return cE
}

func unregister() error {
	return state.c.Close()
}

func systems() ([]lsm.System, error) {
	return state.c.Systems()
}

func pools(search ...string) ([]lsm.Pool, error) {
	if len(search) > 0 {
		return state.c.Pools(search[0], search[1])
	}
	return state.c.Pools()
}

func volumes(search ...string) ([]lsm.Volume, error) {
	if len(search) > 0 {
		return state.c.Volumes(search[0], search[1])
	}
	return state.c.Volumes()
}

func tmoSet(timeout uint32) error {
	state.tmo = timeout
	return nil
}

func tmoGet() uint32 {
	return state.tmo
}

func capabilities(system *lsm.System) (*lsm.Capabilities, error) {
	return state.c.Capabilities(system)
}

func jobStatus(jobID string) (*lsm.JobInfo, error) {
	var item interface{}
	jobStatus, jobPercent, error := state.c.JobStatus(jobID, &item)
	if error != nil {
		return nil, error
	}
	return &lsm.JobInfo{Status: jobStatus, Percent: jobPercent, Item: item}, nil
}

func jobFree(jobID string) error {
	return state.c.JobFree(jobID)
}

func volCreate(pool *lsm.Pool, volumeName string, size uint64,
	provisioning lsm.VolumeProvisionType) (*lsm.Volume, *string, error) {
	return state.c.VolumeCreate(pool, volumeName, size, provisioning, false)
}

func volDelete(volume *lsm.Volume) (*string, error) {
	return state.c.VolumeDelete(volume, false)
}

func disks() ([]lsm.Disk, error) {
	return state.c.Disks()
}

func volReplicate(optionalPool *lsm.Pool, repType lsm.VolumeReplicateType,
	sourceVolume *lsm.Volume, name string) (*lsm.Volume, *string, error) {
	return state.c.VolumeReplicate(optionalPool, repType, sourceVolume, name, false)
}

func volReplicateRange(repType lsm.VolumeReplicateType, srcVol *lsm.Volume, dstVol *lsm.Volume,
	ranges []lsm.BlockRange) (*string, error) {
	return state.c.VolumeReplicateRange(repType, srcVol, dstVol, ranges, false)
}

func volRepRangeBlockSize(system *lsm.System) (uint32, error) {
	return state.c.VolumeRepRangeBlkSize(system)
}

func volResize(vol *lsm.Volume, newSizeBytes uint64) (*lsm.Volume, *string, error) {
	return state.c.VolumeResize(vol, newSizeBytes, false)
}

func volEnable(vol *lsm.Volume) error {
	return state.c.VolumeEnable(vol)
}

func volDisable(vol *lsm.Volume) error {
	return state.c.VolumeDisable(vol)
}

func accessGroups() ([]lsm.AccessGroup, error) {
	return state.c.AccessGroups()
}

func accessGroupCreate(name string, initID string,
	initType lsm.InitiatorType, system *lsm.System) (*lsm.AccessGroup, error) {
	return state.c.AccessGroupCreate(name, initID, initType, system)
}

func accessGroupDelete(ag *lsm.AccessGroup) error {
	return state.c.AccessGroupDelete(ag)
}

func accessGroupInitAdd(ag *lsm.AccessGroup,
	initID string, initType lsm.InitiatorType) (*lsm.AccessGroup, error) {
	return state.c.AccessGroupInitAdd(ag, initID, initType)
}

func accessGroupInitDelete(ag *lsm.AccessGroup,
	initID string, initType lsm.InitiatorType) (*lsm.AccessGroup, error) {
	return state.c.AccessGroupInitDelete(ag, initID, initType)
}

func volumeMask(vol *lsm.Volume, ag *lsm.AccessGroup) error {
	return state.c.VolumeMask(vol, ag)
}

func volumeUnMask(vol *lsm.Volume, ag *lsm.AccessGroup) error {
	return state.c.VolumeUnMask(vol, ag)
}

func volsMaskedToAg(ag *lsm.AccessGroup) ([]lsm.Volume, error) {
	return state.c.VolsMaskedToAg(ag)
}

func agsGrantedToVol(vol *lsm.Volume) ([]lsm.AccessGroup, error) {
	return state.c.AgsGrantedToVol(vol)
}

func iscsiChapAuthSet(initID string, inUser *string, inPassword *string, outUser *string, outPassword *string) error {
	return state.c.IscsiChapAuthSet(initID, inUser, inPassword, outUser, outPassword)
}

func volHasChildDep(vol *lsm.Volume) (bool, error) {
	return state.c.VolHasChildDep(vol)
}

func volChildDepRm(vol *lsm.Volume) (*string, error) {
	return state.c.VolChildDepRm(vol, false)
}

func targetPorts() ([]lsm.TargetPort, error) {
	return state.c.TargetPorts()
}

func volIdentLedOn(vol *lsm.Volume) error {
	return state.c.VolIdentLedOn(vol)
}

func volIdentLedOff(vol *lsm.Volume) error {
	return state.c.VolIdentLedOff(vol)
}

func fileSystems(search ...string) ([]lsm.FileSystem, error) {
	if len(search) > 0 {
		return state.c.FileSystems(search[0], search[1])
	}
	return state.c.FileSystems()
}

func fileSystemCreate(pool *lsm.Pool, name string, sizeBytes uint64) (*lsm.FileSystem, *string, error) {
	return state.c.FsCreate(pool, name, sizeBytes, false)
}

func fileSystemDelete(fs *lsm.FileSystem) (*string, error) {
	return state.c.FsDelete(fs, false)
}

func fileSystemResize(fs *lsm.FileSystem, newSizeBytes uint64) (*lsm.FileSystem, *string, error) {
	return state.c.FsResize(fs, newSizeBytes, false)
}

func fileSystemClone(srcFs *lsm.FileSystem, destName string, optionalSnapShot *lsm.FileSystemSnapShot) (*lsm.FileSystem, *string, error) {
	return state.c.FsClone(srcFs, destName, optionalSnapShot, false)
}

func fileSystemFileClone(fs *lsm.FileSystem, srcFileName string, dstFileName string,
	optionalSnapShot *lsm.FileSystemSnapShot) (*string, error) {
	return state.c.FsFileClone(fs, srcFileName, dstFileName, optionalSnapShot, false)
}

func fileSystemSnapShotCreate(fs *lsm.FileSystem, name string) (*lsm.FileSystemSnapShot, *string, error) {
	return state.c.FsSnapShotCreate(fs, name, false)
}

func fileSystemSnapShotDelete(fs *lsm.FileSystem, ss *lsm.FileSystemSnapShot) (*string, error) {
	return state.c.FsSnapShotDelete(fs, ss, false)
}

func fileSystemSnapShots(fs *lsm.FileSystem) ([]lsm.FileSystemSnapShot, error) {
	return state.c.FsSnapShots(fs)
}

func fileSystemSnapShotRestore(fs *lsm.FileSystem, ss *lsm.FileSystemSnapShot,
	allFiles bool, files []string, restoreFiles []string) (*string, error) {

	return state.c.FsSnapShotRestore(fs, ss, allFiles, files, restoreFiles, false)
}

func fileSystemHasChildDep(fs *lsm.FileSystem, files []string) (bool, error) {
	return state.c.FsHasChildDep(fs, files)
}

func fileSystemChildDepRm(fs *lsm.FileSystem, files []string) (*string, error) {
	return state.c.FsChildDepRm(fs, files, false)
}

func exports(search ...string) ([]lsm.NfsExport, error) {
	if len(search) > 1 {
		return state.c.NfsExports(search[0], search[1])
	}
	return state.c.NfsExports()
}

func fsExport(fs *lsm.FileSystem, exportPath *string,
	access *lsm.NfsAccess, authType *string, options *string) (*lsm.NfsExport, error) {
	return state.c.FsExport(fs, exportPath, access, authType, options)
}

func fsUnExport(export *lsm.NfsExport) error {
	return state.c.FsUnExport(export)
}

func fsExportAuthTypes() ([]string, error) {
	return state.c.NfsExportAuthTypes()
}

func volRaidCreate(name string, raidType lsm.RaidType, disks []lsm.Disk, stripSize uint32) (*lsm.Volume, error) {
	return state.c.VolRaidCreate(name, raidType, disks, stripSize)
}

func volRaidCreateCapGet(sys *lsm.System) (*lsm.SupportedRaidCapability, error) {
	return state.c.VolRaidCreateCapGet(sys)
}

func poolMemberInfo(pool *lsm.Pool) (*lsm.PoolMemberInfo, error) {
	return state.c.PoolMemberInfo(pool)
}

func handleVolRaidInfo(vol *lsm.Volume) (*lsm.VolumeRaidInfo, error) {
	return state.c.VolRaidInfo(vol)
}

func handleBatteries() ([]lsm.Battery, error) {
	return state.c.Batteries()
}

func main() {
	var cb lsm.PluginCallBacks
	cb.Mgmt.Systems = systems
	cb.Mgmt.PluginRegister = register
	cb.Mgmt.PluginUnregister = unregister
	cb.Mgmt.Pools = pools
	cb.Mgmt.TimeOutSet = tmoSet
	cb.Mgmt.TimeOutGet = tmoGet
	cb.Mgmt.Capabilities = capabilities
	cb.Mgmt.JobStatus = jobStatus
	cb.Mgmt.JobFree = jobFree

	cb.San.VolumeCreate = volCreate
	cb.San.VolumeDelete = volDelete
	cb.San.Volumes = volumes
	cb.San.Disks = disks
	cb.San.VolumeReplicate = volReplicate
	cb.San.VolumeReplicateRange = volReplicateRange
	cb.San.VolumeRepRangeBlkSize = volRepRangeBlockSize
	cb.San.VolumeResize = volResize
	cb.San.VolumeEnable = volEnable
	cb.San.VolumeDisable = volDisable
	cb.San.VolumeMask = volumeMask
	cb.San.VolumeUnMask = volumeUnMask
	cb.San.VolsMaskedToAg = volsMaskedToAg
	cb.San.VolHasChildDep = volHasChildDep
	cb.San.VolChildDepRm = volChildDepRm

	cb.San.AccessGroups = accessGroups
	cb.San.AccessGroupCreate = accessGroupCreate
	cb.San.AccessGroupDelete = accessGroupDelete
	cb.San.AccessGroupInitAdd = accessGroupInitAdd
	cb.San.AccessGroupInitDelete = accessGroupInitDelete
	cb.San.AgsGrantedToVol = agsGrantedToVol
	cb.San.IscsiChapAuthSet = iscsiChapAuthSet

	cb.San.TargetPorts = targetPorts
	cb.San.VolIdentLedOn = volIdentLedOn
	cb.San.VolIdentLedOff = volIdentLedOff

	cb.File.FileSystems = fileSystems
	cb.File.FsCreate = fileSystemCreate
	cb.File.FsDelete = fileSystemDelete
	cb.File.FsResize = fileSystemResize
	cb.File.FsClone = fileSystemClone
	cb.File.FsFileClone = fileSystemFileClone
	cb.File.FsSnapShotCreate = fileSystemSnapShotCreate
	cb.File.FsSnapShotDelete = fileSystemSnapShotDelete
	cb.File.FsSnapShots = fileSystemSnapShots
	cb.File.FsSnapShotRestore = fileSystemSnapShotRestore

	cb.File.FsHasChildDep = fileSystemHasChildDep
	cb.File.FsChildDepRm = fileSystemChildDepRm

	cb.Nfs.Exports = exports
	cb.Nfs.FsExport = fsExport
	cb.Nfs.FsUnExport = fsUnExport
	cb.Nfs.ExportAuthTypes = fsExportAuthTypes

	cb.Hba.VolRaidCreate = volRaidCreate
	cb.Hba.VolRaidCreateCapGet = volRaidCreateCapGet
	cb.Hba.PoolMemberInfo = poolMemberInfo
	cb.Hba.VolRaidInfo = handleVolRaidInfo
	cb.Hba.Batteries = handleBatteries
	plugin, err := lsm.PluginInit(&cb, os.Args, "golang forwarding plugin", "0.0.1")
	if err != nil {
		fmt.Printf("Failed to initialize plugin, exiting! (%s)\n", err)
	} else {
		plugin.Run()
	}
}
