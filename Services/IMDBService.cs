using System.Net.Http;
using System.Threading.Tasks;

namespace Finderr.Services
{
    public class IMDBService
    {
        private readonly HttpClient _httpClient;
        private readonly string _apiKey = "imdb_api_key";

        public IMDBService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<string> GetMovieMetadata(string imdbId)
        {
            var response = await _httpClient.GetAsync($"https://imdb-api.com/en/API/Title/{_apiKey}/{imdbId}");
            return await response.Content.ReadAsStringAsync();
        }
    }
}
