// Copyright 2023 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package linux

// Constants used by keyctl(2) and other keyrings-related syscalls.
// Source: include/uapi/linux/keyctl.h

const (
	KEY_SPEC_THREAD_KEYRING       = -1
	KEY_SPEC_PROCESS_KEYRING      = -2
	KEY_SPEC_SESSION_KEYRING      = -3
	KEY_SPEC_USER_KEYRING         = -4
	KEY_SPEC_USER_SESSION_KEYRING = -5
	KEY_SPEC_GROUP_KEYRING        = -6
	KEY_SPEC_REQKEY_AUTH_KEY      = -7
	KEY_SPEC_REQUESTOR_KEYRING    = -8
)

const (
	KEY_REQKEY_DEFL_NO_CHANGE            = -1
	KEY_REQKEY_DEFL_DEFAULT              = 0
	KEY_REQKEY_DEFL_THREAD_KEYRING       = 1
	KEY_REQKEY_DEFL_PROCESS_KEYRING      = 2
	KEY_REQKEY_DEFL_SESSION_KEYRING      = 3
	KEY_REQKEY_DEFL_USER_KEYRING         = 4
	KEY_REQKEY_DEFL_USER_SESSION_KEYRING = 5
	KEY_REQKEY_DEFL_GROUP_KEYRING        = 6
	KEY_REQKEY_DEFL_REQUESTOR_KEYRING    = 7
)

const (
	KEYCTL_GET_KEYRING_ID       = 0
	KEYCTL_JOIN_SESSION_KEYRING = 1
	KEYCTL_UPDATE               = 2
	KEYCTL_REVOKE               = 3
	KEYCTL_CHOWN                = 4
	KEYCTL_SETPERM              = 5
	KEYCTL_DESCRIBE             = 6
	KEYCTL_CLEAR                = 7
	KEYCTL_LINK                 = 8
	KEYCTL_UNLINK               = 9
	KEYCTL_SEARCH               = 10
	KEYCTL_READ                 = 11
	KEYCTL_INSTANTIATE          = 12
	KEYCTL_NEGATE               = 13
	KEYCTL_SET_REQKEY_KEYRING   = 14
	KEYCTL_SET_TIMEOUT          = 15
	KEYCTL_ASSUME_AUTHORITY     = 16
	KEYCTL_GET_SECURITY         = 17
	KEYCTL_SESSION_TO_PARENT    = 18
	KEYCTL_REJECT               = 19
	KEYCTL_INSTANTIATE_IOV      = 20
	KEYCTL_INVALIDATE           = 21
	KEYCTL_GET_PERSISTENT       = 22
	KEYCTL_DH_COMPUTE           = 23
	KEYCTL_PKEY_QUERY           = 24
	KEYCTL_PKEY_ENCRYPT         = 25
	KEYCTL_PKEY_DECRYPT         = 26
	KEYCTL_PKEY_SIGN            = 27
	KEYCTL_PKEY_VERIFY          = 28
	KEYCTL_RESTRICT_KEYRING     = 29
	KEYCTL_MOVE                 = 30
	KEYCTL_CAPABILITIES         = 31
	KEYCTL_WATCH_KEY            = 32
)

const (
	KEYCTL_SUPPORTS_ENCRYPT = 0x01
	KEYCTL_SUPPORTS_DECRYPT = 0x02
	KEYCTL_SUPPORTS_SIGN    = 0x04
	KEYCTL_SUPPORTS_VERIFY  = 0x08
)

const (
	KEYCTL_CAPS0_CAPABILITIES        = 0x01
	KEYCTL_CAPS0_PERSISTENT_KEYRINGS = 0x02
	KEYCTL_CAPS0_DIFFIE_HELLMAN      = 0x04
	KEYCTL_CAPS0_PUBLIC_KEY          = 0x08
	KEYCTL_CAPS0_BIG_KEY             = 0x10
	KEYCTL_CAPS0_INVALIDATE          = 0x20
	KEYCTL_CAPS0_RESTRICT_KEYRING    = 0x40
	KEYCTL_CAPS0_MOVE                = 0x80
	KEYCTL_CAPS1_NS_KEYRING_NAME     = 0x01
	KEYCTL_CAPS1_NS_KEY_TAG          = 0x02
	KEYCTL_CAPS1_NOTIFICATIONS       = 0x04
)
