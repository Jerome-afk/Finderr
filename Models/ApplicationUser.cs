using Microsoft.AspNetCore.Identity;

namespace Finderr.Models
{
    public class ApplicationUser : IdentityUser
    {
        public string? FirstName { get; set; }
        public string? LastName { get; set; }
        //public string? UserName { get; set; }
        public DateTime DateRegistered { get; set; } = DateTime.UtcNow;
    }
}
