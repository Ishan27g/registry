package mapper

import domain.PeerJson
import domain.Peer
import domain.Zone

object RegistryPeers {
    private val numZones: String = (System.getenv("Zones") ?: "3")
    private var zones= mutableMapOf<Int, Zone>()
    private var currentZone: Int = 1
    init {
        for (zoneId in 1..numZones.toInt()){
            zones[zoneId] = Zone(zoneId)
        }
    }
    private fun nextZone(): Int{
        val zone = currentZone
        currentZone++
        if (currentZone > numZones.toInt()){
            currentZone = 1;
        }
        return zone
    }
    fun toJson():Map<Int,List<PeerJson>>{
        synchronized(this){
            val data= mutableMapOf<Int, List<PeerJson>>()
            zones.forEach{
                if (it.value.peers.isNotEmpty()){
                    println(it.value.peers.size)
                    data[it.key] = it.value.toJson
                }
            }
            return data
        }
    }
    fun register(address: String, metadata: Map<String, Any>):Int {
        synchronized(this){
            val zoneId = nextZone()
            val currentZonePeers = zones[zoneId]
            return if (currentZonePeers != null) {
                var peer = Peer(address, metadata, zoneId)
                currentZonePeers.add(peer)
                zones[zoneId] = currentZonePeers
                zoneId
            }else{
                400
            }
        }
    }
    fun get(address: String): PeerJson? {
        synchronized(this){
            for (zone in zones){
                zone.value.peers.forEach {
                    if (it.value.toJson().address.compareTo(address) == 0){
                        return it.value.toJson()
                    }
                }
            }
        }
        return null
    }
    fun get(zoneId: Int): List<PeerJson>?{
        synchronized(this) {
            return zones[zoneId]?.toJson
        }
    }
    fun getZoneIds(): List<Int>{
        synchronized(this) {
            val rsp = mutableListOf<Int>()
            for (zones in zones.values) {
                if (zones.peers.isNotEmpty()) {
                    rsp.add(zones.zoneId)
                }
            }
            return rsp
        }
    }

    fun remove(address: String){
        synchronized(this) {
            if (get(address) !=null){
                for (zone in zones) {
                    zone.value.peers = zone.value.peers.filter {
                        it.value.toJson().address != address
                    } as MutableMap<String, Peer>
                }
            }
        }
    }
}
