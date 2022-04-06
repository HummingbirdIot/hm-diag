#!/bin/sh

CONTAINER_NAMES=$(docker ps | grep -oE '[a-z0-9_-]+helium-miner_1')

docker exec $CONTAINER_NAMES miner eval "\
Ledger = blockchain:ledger(blockchain_worker:blockchain()),\
PubKey = blockchain_swarm:pubkey_bin(),\
case blockchain_ledger_v1:find_gateway_info(PubKey, Ledger) of\
    {ok, Gateway} ->\
         GWLoc = case blockchain_ledger_gateway_v2:location(Gateway) of \
                    undefined -> 'none'; \
                    L -> h3:to_string(L) \
                    end, \
         GWOwnAddr = libp2p_crypto:pubkey_bin_to_p2p(blockchain_ledger_gateway_v2:owner_address(Gateway)),\
         [lists:concat(['location: ', GWLoc]), lists:concat(['owner: ', GWOwnAddr])];\
    _ ->\
         [lists:concat(['location: ', 'none']), lists:concat(['owner: ', '/p2p/none'])]\
end."
