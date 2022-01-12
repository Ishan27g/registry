package domain

import mapper.RegistryPeers
import org.testng.annotations.Test

import org.testng.Assert.*

class RegistryPeersTest {
    private fun mockPeer(address:String): PeerJson {
        return Peer(address, mutableMapOf<String,Any>(), 1).toJson()
    }
    @Test
    fun testToJson() {
        val mp = mockPeer("1")
        assertEquals("1", mp.address)
        assertEquals(1, mp.zone)
    }

    @Test
    fun testRegister() {
        var pm = RegistryPeers
        val mp = mockPeer("1")
        pm.register(mp.address, mp.meta_data)
        assertNotNull(pm.get("1"))
    }

    @Test
    fun testGet() {
        var pm = RegistryPeers
        val mp1 = mockPeer("1")
        val mp2 = mockPeer("2")
        pm.register(mp1.address, mp1.meta_data)
        pm.register(mp2.address, mp2.meta_data)
        assertNotNull(pm.get("1"))
        assertNotNull(pm.get("2"))
    }

    @Test
    fun testRemove() {
        var pm = RegistryPeers
        val mp = mockPeer("1")
        pm.register(mp.address, mp.meta_data)
        pm.remove("1")
        assertNull(pm.get("1"))
    }
}