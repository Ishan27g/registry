package domain

import java.net.http.HttpClient
import kotlinx.coroutines.*
import java.lang.Exception
import java.net.URI
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import java.time.Duration

@ExperimentalCoroutinesApi
object Monitor {
    private val client = HttpClient.newBuilder().build()

    private suspend fun monitor(address: String) {
        val request = HttpRequest.newBuilder()
            .uri(URI.create(address))
            .timeout(Duration.ofSeconds(3))
            .build();
        var retries = 0
        while (true) {
            delay(3000)
            val response: HttpResponse<String>
            try {
                response = client.send(request, HttpResponse.BodyHandlers.ofString())
                if (response.statusCode() != 200){
                    retries ++
                    if (retries == 2){
                        break
                    }
                }
            }catch (e : Exception){
                break
            }
        }
    }
    // monitorPeer launches a coroutine that runs as long as peer is reachable
    suspend fun monitorPeer(address: String){
        coroutineScope {
            launch {
                monitor(address)
            }
        }
    }

}