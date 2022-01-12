package domain

import java.time.Instant

internal class Peer(val address: String, private val metadata: Map<String, Any>, private val zoneId: Int) {
    private var registeredAt: Instant = Instant.now()

    fun toJson(): PeerJson {
        return PeerJson(this.address, this.metadata, this.registeredAt.toString(), this.zoneId)
    }

    fun registeredBefore(peer: Peer):Boolean{
        return this.registeredAt.isBefore(peer.registeredAt)
    }
}