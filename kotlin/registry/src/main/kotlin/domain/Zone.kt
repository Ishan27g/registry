package domain

internal class Zone(val zoneId: Int){
    var peers = mutableMapOf<String, Peer>()
    fun add(peer: Peer){
        if (peers[peer.address] == null){
            peers[peer.address] = peer
            return
        }
        if (peers[peer.address]!!.registeredBefore(peer)){
            peers[peer.address] = peer
            return
        }
    }

    val toJson: List<PeerJson>
        get() {
            val rsp = mutableListOf<PeerJson>()
            peers.forEach {
                rsp.add(it.value.toJson())
            }
            return rsp
        }
}