import controller.Registry
import io.javalin.Javalin
import kotlinx.coroutines.ExperimentalCoroutinesApi

@ExperimentalCoroutinesApi
class Server(){
    private val app:Javalin = Javalin.create{ config ->
        config.requestLogger { ctx,_ ->
            println( "["+ctx.method()+"]" + " " + ctx.url()) }
    }
    private val registry = Registry // static singleton
    init {
        app.get("/reset") { ctx -> registry.reset(ctx) }
        app.get("/zones") { ctx -> registry.getZonesIds(ctx) }
        app.get("/details") { ctx -> registry.getPeers(ctx) }
        app.get("/zone/{id}") { ctx -> registry.getPeersForZone(ctx) }
        app.post("/register") { ctx -> registry.register(ctx) }
    }
    fun start(){
        val port: String = (System.getenv("Port") ?: "9999")
        app.start(port.toInt())
    }
}

fun main(){
    Server().start()
}