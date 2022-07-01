#!/bin/bash
# For details about this script either read the README.md file or the
# relevant comment section of the python script find_candidate_same_files.py
mkdir -p georzaza_results
mkdir -p georzaza_results/ethereum
awk 'NR==FNR{arr[$0];next} $0 in arr' ./ethereum.go ./types/ethereum.go > georzaza_results/ethereum/types_ethereum
mkdir -p georzaza_results/natpmp
awk 'NR==FNR{arr[$0];next} $0 in arr' ./natpmp.go ./p2p/natpmp.go > georzaza_results/natpmp/p2p_natpmp
mkdir -p georzaza_results/events
awk 'NR==FNR{arr[$0];next} $0 in arr' ./events.go ./core/events.go > georzaza_results/events/core_events
mkdir -p georzaza_results/natupnp
awk 'NR==FNR{arr[$0];next} $0 in arr' ./natupnp.go ./p2p/natupnp.go > georzaza_results/natupnp/p2p_natupnp
mkdir -p georzaza_results/peer
awk 'NR==FNR{arr[$0];next} $0 in arr' ./peer.go ./p2p/peer.go > georzaza_results/peer/p2p_peer
awk 'NR==FNR{arr[$0];next} $0 in arr' ./peer.go ./whisper/peer.go > georzaza_results/peer/whisper_peer
mkdir -p georzaza_results/p2p_peer
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/peer.go ./whisper/peer.go > georzaza_results/p2p_peer/whisper_peer
mkdir -p georzaza_results/trie_iterator
awk 'NR==FNR{arr[$0];next} $0 in arr' ./trie/iterator.go ./ptrie/iterator.go > georzaza_results/trie_iterator/ptrie_iterator
mkdir -p georzaza_results/trie_trie_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./trie/trie_test.go ./ptrie/trie_test.go > georzaza_results/trie_trie_test/ptrie_trie_test
mkdir -p georzaza_results/trie_main_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./trie/main_test.go ./state/main_test.go > georzaza_results/trie_main_test/state_main_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./trie/main_test.go ./vm/main_test.go > georzaza_results/trie_main_test/vm_main_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./trie/main_test.go ./ethutil/main_test.go > georzaza_results/trie_main_test/ethutil_main_test
mkdir -p georzaza_results/state_main_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./state/main_test.go ./vm/main_test.go > georzaza_results/state_main_test/vm_main_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./state/main_test.go ./ethutil/main_test.go > georzaza_results/state_main_test/ethutil_main_test
mkdir -p georzaza_results/vm_main_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./vm/main_test.go ./ethutil/main_test.go > georzaza_results/vm_main_test/ethutil_main_test
mkdir -p georzaza_results/trie_trie
awk 'NR==FNR{arr[$0];next} $0 in arr' ./trie/trie.go ./ptrie/trie.go > georzaza_results/trie_trie/ptrie_trie
awk 'NR==FNR{arr[$0];next} $0 in arr' ./trie/trie.go ./tests/helper/trie.go > georzaza_results/trie_trie/tests_helper_trie
mkdir -p georzaza_results/ptrie_trie
awk 'NR==FNR{arr[$0];next} $0 in arr' ./ptrie/trie.go ./tests/helper/trie.go > georzaza_results/ptrie_trie/tests_helper_trie
mkdir -p georzaza_results/p2p_server
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/server.go ./websocket/server.go > georzaza_results/p2p_server/websocket_server
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/server.go ./rpc/server.go > georzaza_results/p2p_server/rpc_server
mkdir -p georzaza_results/websocket_server
awk 'NR==FNR{arr[$0];next} $0 in arr' ./websocket/server.go ./rpc/server.go > georzaza_results/websocket_server/rpc_server
mkdir -p georzaza_results/p2p_client_identity_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/client_identity_test.go ./wire/client_identity_test.go > georzaza_results/p2p_client_identity_test/wire_client_identity_test
mkdir -p georzaza_results/p2p_message
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/message.go ./websocket/message.go > georzaza_results/p2p_message/websocket_message
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/message.go ./whisper/message.go > georzaza_results/p2p_message/whisper_message
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/message.go ./rpc/message.go > georzaza_results/p2p_message/rpc_message
mkdir -p georzaza_results/websocket_message
awk 'NR==FNR{arr[$0];next} $0 in arr' ./websocket/message.go ./whisper/message.go > georzaza_results/websocket_message/whisper_message
awk 'NR==FNR{arr[$0];next} $0 in arr' ./websocket/message.go ./rpc/message.go > georzaza_results/websocket_message/rpc_message
mkdir -p georzaza_results/whisper_message
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/message.go ./rpc/message.go > georzaza_results/whisper_message/rpc_message
mkdir -p georzaza_results/p2p_client_identity
awk 'NR==FNR{arr[$0];next} $0 in arr' ./p2p/client_identity.go ./wire/client_identity.go > georzaza_results/p2p_client_identity/wire_client_identity
mkdir -p georzaza_results/core_asm
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/asm.go ./vm/asm.go > georzaza_results/core_asm/vm_asm
mkdir -p georzaza_results/core_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/filter.go ./whisper/filter.go > georzaza_results/core_filter/whisper_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/filter.go ./event/filter/filter.go > georzaza_results/core_filter/event_filter_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/filter.go ./ui/filter.go > georzaza_results/core_filter/ui_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/filter.go ./ui/qt/filter.go > georzaza_results/core_filter/ui_qt_filter
mkdir -p georzaza_results/whisper_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/filter.go ./event/filter/filter.go > georzaza_results/whisper_filter/event_filter_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/filter.go ./ui/filter.go > georzaza_results/whisper_filter/ui_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/filter.go ./ui/qt/filter.go > georzaza_results/whisper_filter/ui_qt_filter
mkdir -p georzaza_results/event_filter_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./event/filter/filter.go ./ui/filter.go > georzaza_results/event_filter_filter/ui_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./event/filter/filter.go ./ui/qt/filter.go > georzaza_results/event_filter_filter/ui_qt_filter
mkdir -p georzaza_results/ui_filter
awk 'NR==FNR{arr[$0];next} $0 in arr' ./ui/filter.go ./ui/qt/filter.go > georzaza_results/ui_filter/ui_qt_filter
mkdir -p georzaza_results/core_vm_env
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/vm_env.go ./xeth/vm_env.go > georzaza_results/core_vm_env/xeth_vm_env
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/vm_env.go ./cmd/utils/vm_env.go > georzaza_results/core_vm_env/cmd_utils_vm_env
mkdir -p georzaza_results/xeth_vm_env
awk 'NR==FNR{arr[$0];next} $0 in arr' ./xeth/vm_env.go ./cmd/utils/vm_env.go > georzaza_results/xeth_vm_env/cmd_utils_vm_env
mkdir -p georzaza_results/core_filter_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/filter_test.go ./event/filter/filter_test.go > georzaza_results/core_filter_test/event_filter_filter_test
mkdir -p georzaza_results/core_types_block
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/types/block.go ./pow/block.go > georzaza_results/core_types_block/pow_block
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/types/block.go ./pow/ar/block.go > georzaza_results/core_types_block/pow_ar_block
mkdir -p georzaza_results/pow_block
awk 'NR==FNR{arr[$0];next} $0 in arr' ./pow/block.go ./pow/ar/block.go > georzaza_results/pow_block/pow_ar_block
mkdir -p georzaza_results/core_types_common
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/types/common.go ./vm/common.go > georzaza_results/core_types_common/vm_common
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/types/common.go ./ethutil/common.go > georzaza_results/core_types_common/ethutil_common
awk 'NR==FNR{arr[$0];next} $0 in arr' ./core/types/common.go ./tests/helper/common.go > georzaza_results/core_types_common/tests_helper_common
mkdir -p georzaza_results/vm_common
awk 'NR==FNR{arr[$0];next} $0 in arr' ./vm/common.go ./ethutil/common.go > georzaza_results/vm_common/ethutil_common
awk 'NR==FNR{arr[$0];next} $0 in arr' ./vm/common.go ./tests/helper/common.go > georzaza_results/vm_common/tests_helper_common
mkdir -p georzaza_results/ethutil_common
awk 'NR==FNR{arr[$0];next} $0 in arr' ./ethutil/common.go ./tests/helper/common.go > georzaza_results/ethutil_common/tests_helper_common
mkdir -p georzaza_results/state_errors
awk 'NR==FNR{arr[$0];next} $0 in arr' ./state/errors.go ./vm/errors.go > georzaza_results/state_errors/vm_errors
awk 'NR==FNR{arr[$0];next} $0 in arr' ./state/errors.go ./cmd/mist/errors.go > georzaza_results/state_errors/cmd_mist_errors
mkdir -p georzaza_results/vm_errors
awk 'NR==FNR{arr[$0];next} $0 in arr' ./vm/errors.go ./cmd/mist/errors.go > georzaza_results/vm_errors/cmd_mist_errors
mkdir -p georzaza_results/javascript_types
awk 'NR==FNR{arr[$0];next} $0 in arr' ./javascript/types.go ./vm/types.go > georzaza_results/javascript_types/vm_types
mkdir -p georzaza_results/whisper_whisper
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/whisper.go ./ui/qt/qwhisper/whisper.go > georzaza_results/whisper_whisper/ui_qt_qwhisper_whisper
mkdir -p georzaza_results/whisper_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/main.go ./cmd/ethereum/main.go > georzaza_results/whisper_main/cmd_ethereum_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/main.go ./cmd/ethtest/main.go > georzaza_results/whisper_main/cmd_ethtest_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/main.go ./cmd/peerserver/main.go > georzaza_results/whisper_main/cmd_peerserver_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/main.go ./cmd/mist/main.go > georzaza_results/whisper_main/cmd_mist_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./whisper/main.go ./cmd/evm/main.go > georzaza_results/whisper_main/cmd_evm_main
mkdir -p georzaza_results/cmd_ethereum_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethereum/main.go ./cmd/ethtest/main.go > georzaza_results/cmd_ethereum_main/cmd_ethtest_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethereum/main.go ./cmd/peerserver/main.go > georzaza_results/cmd_ethereum_main/cmd_peerserver_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethereum/main.go ./cmd/mist/main.go > georzaza_results/cmd_ethereum_main/cmd_mist_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethereum/main.go ./cmd/evm/main.go > georzaza_results/cmd_ethereum_main/cmd_evm_main
mkdir -p georzaza_results/cmd_ethtest_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethtest/main.go ./cmd/peerserver/main.go > georzaza_results/cmd_ethtest_main/cmd_peerserver_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethtest/main.go ./cmd/mist/main.go > georzaza_results/cmd_ethtest_main/cmd_mist_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethtest/main.go ./cmd/evm/main.go > georzaza_results/cmd_ethtest_main/cmd_evm_main
mkdir -p georzaza_results/cmd_peerserver_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/peerserver/main.go ./cmd/mist/main.go > georzaza_results/cmd_peerserver_main/cmd_mist_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/peerserver/main.go ./cmd/evm/main.go > georzaza_results/cmd_peerserver_main/cmd_evm_main
mkdir -p georzaza_results/cmd_mist_main
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/mist/main.go ./cmd/evm/main.go > georzaza_results/cmd_mist_main/cmd_evm_main
mkdir -p georzaza_results/vm_vm
awk 'NR==FNR{arr[$0];next} $0 in arr' ./vm/vm.go ./tests/helper/vm.go > georzaza_results/vm_vm/tests_helper_vm
mkdir -p georzaza_results/vm_debugger
awk 'NR==FNR{arr[$0];next} $0 in arr' ./vm/debugger.go ./cmd/mist/debugger.go > georzaza_results/vm_debugger/cmd_mist_debugger
mkdir -p georzaza_results/pow_pow
awk 'NR==FNR{arr[$0];next} $0 in arr' ./pow/pow.go ./pow/ezp/pow.go > georzaza_results/pow_pow/pow_ezp_pow
awk 'NR==FNR{arr[$0];next} $0 in arr' ./pow/pow.go ./pow/ar/pow.go > georzaza_results/pow_pow/pow_ar_pow
mkdir -p georzaza_results/pow_ezp_pow
awk 'NR==FNR{arr[$0];next} $0 in arr' ./pow/ezp/pow.go ./pow/ar/pow.go > georzaza_results/pow_ezp_pow/pow_ar_pow
mkdir -p georzaza_results/logger_example_test
awk 'NR==FNR{arr[$0];next} $0 in arr' ./logger/example_test.go ./event/example_test.go > georzaza_results/logger_example_test/event_example_test
mkdir -p georzaza_results/xeth_config
awk 'NR==FNR{arr[$0];next} $0 in arr' ./xeth/config.go ./ethutil/config.go > georzaza_results/xeth_config/ethutil_config
mkdir -p georzaza_results/cmd_utils_cmd
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/utils/cmd.go ./cmd/ethereum/cmd.go > georzaza_results/cmd_utils_cmd/cmd_ethereum_cmd
mkdir -p georzaza_results/cmd_ethereum_flags
awk 'NR==FNR{arr[$0];next} $0 in arr' ./cmd/ethereum/flags.go ./cmd/mist/flags.go > georzaza_results/cmd_ethereum_flags/cmd_mist_flags
