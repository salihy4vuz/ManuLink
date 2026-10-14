using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace SheetMetalApex.Bridge
{
    public class WebSocketClient
    {
        private ClientWebSocket _client;
        private Uri _uri;

        public WebSocketClient(string url)
        {
            _uri = new Uri(url);
            _client = new ClientWebSocket();
        }

        public WebSocketClient()
        {
            _uri = new Uri(GetWebSocketUrl());
            _client = new ClientWebSocket();
        }

        private static string GetWebSocketUrl()
        {
            bool useTLS = Environment.GetEnvironmentVariable("SMA_USE_TLS") == "true";
            string host = Environment.GetEnvironmentVariable("SMA_API_HOST") ?? "localhost:3000";
            return useTLS ? $"wss://{host}/ws" : $"ws://{host}/ws";
        }

        public async Task ConnectAsync()
        {
            while (true)
            {
                try
                {
                    if (_client.State != WebSocketState.Open)
                    {
                        if (_client.State == WebSocketState.Aborted || _client.State == WebSocketState.Closed)
                        {
                            _client.Dispose();
                            _client = new ClientWebSocket();
                        }
                        await _client.ConnectAsync(_uri, CancellationToken.None);
                        Console.WriteLine("Golang B2B API (WebSocket) bağlantısı başarılı.");
                    }
                    await Task.Delay(5000); // Heartbeat
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Bağlantı koptu, yeniden deneniyor: {ex.Message}");
                    await Task.Delay(5000); // Exponential backoff can be added here
                }
            }
        }

        public async Task SendMessageAsync(string message)
        {
            if (_client.State == WebSocketState.Open)
            {
                byte[] bytes = Encoding.UTF8.GetBytes(message);
                await _client.SendAsync(new ArraySegment<byte>(bytes), WebSocketMessageType.Text, true, CancellationToken.None);
            }
        }

        public async Task DisconnectAsync()
        {
            if (_client.State == WebSocketState.Open)
            {
                await _client.CloseAsync(WebSocketCloseStatus.NormalClosure, "Add-in closed", CancellationToken.None);
            }
            _client.Dispose();
        }
    }
}
