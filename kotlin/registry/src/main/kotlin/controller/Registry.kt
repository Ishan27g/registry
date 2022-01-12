package controller

import domain.Monitor
import domain.PeerJson
import mapper.RegistryPeers

import io.javalin.http.Context
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch


@ExperimentalCoroutinesApi
object Registry {
    private var peerMap = RegistryPeers // static singleton
    private var client = Monitor // static singleton

    private suspend fun monitorPeer(address: String)= coroutineScope {
        println("Monitoring peer - $address")
        client.monitorPeer(address) // blocked until disconnected
        synchronized(this){
            peerMap.remove(address)
        }
        println("Removed inactive peer - $address")
    }
    fun reset(ctx: Context){
        synchronized(this){
            peerMap = RegistryPeers
            client = Monitor
            ctx.result("resetting...")
        }
    }
    fun register(ctx: Context){
        synchronized(this){
            data class JsonReq(val address: String, val metaData: Map<String, Any>)
            val json = ctx.bodyAsClass<JsonReq>()
            if (peerMap.get(json.address) == null){
                GlobalScope.launch { monitorPeer(json.address) }
            }
            val zoneId = peerMap.register(json.address, json.metaData)
            val rsp = peerMap.get(zoneId)
            ctx.json(rsp!!)
        }
    }
    fun getPeersForZone(ctx: Context){
        synchronized(this) {
            val zoneId: String = ctx.pathParam("id")
            val list= listOf<PeerJson>()
            val rsp = peerMap.get(zoneId.toInt())
            ctx.json(rsp ?: list)
        }
    }
    fun getPeers(ctx: Context){
        synchronized(this) {
            ctx.json(peerMap.toJson())
        }
    }
    fun getZonesIds(ctx: Context) {
        synchronized(this) {
            val rsp = peerMap.getZoneIds()
            ctx.json(rsp)
        }
    }
}