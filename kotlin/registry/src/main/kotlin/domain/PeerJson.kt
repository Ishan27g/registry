package domain

data class PeerJson(val address: String, val meta_data: Map<String, Any>,
                    val registered_at: String, var zone:Int)
