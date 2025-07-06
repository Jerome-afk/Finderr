using System.Net.Http;
using System.Threading.Tasks;
using System.Text;
using Newtonsoft.Json;

namespace Finderr.Services
{
    public class AniListService
    {

        private readonly HttpClient _httpClient;

        public AniListService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<string> GetAnimeMetadata(int animeId)
        {
            var query = new
            {
                query = $@"query {{
                    Media(id: {animeId}) {{
                        title {{
                            romanji
                            english
                            native }}
                        description
                        episodes
                        duration
                        genres
                        averageScore
                        coverImage {{ large }}
                    }}
                }}"
            };

            var content = new StringContent(
                JsonConvert.SerializeObject(query),
                Encoding.UTF8,
                "application/json"
            );

            var response = await _httpClient.PostAsync("https://graphql.anilist.co", content);
            if (response.IsSuccessStatusCode)
            {
                var jsonResponse = await response.Content.ReadAsStringAsync();
                return jsonResponse;
            }
            else
            {
                throw new HttpRequestException($"Error fetching data from AniList: {response.ReasonPhrase}");
            }
        }
    }
}
