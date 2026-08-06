using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using EbookStore.Entities.User;
using Microsoft.IdentityModel.Tokens;

public class JwtGenerator
{
    private readonly IConfiguration _configuration;

    public JwtGenerator(IConfiguration configuration)
    {
        _configuration = configuration;
    }

    public string Generate(User user)
    {
        var claims = new[]{
            new Claim(JwtRegisteredClaimNames.Sub, user.Id.ToString()),
            new Claim(JwtRegisteredClaimNames.Email, user.Email),
            new Claim(ClaimTypes.Name,$"{user.FirstName} {user.LastName}"),
            new Claim(ClaimTypes.Role,$"{user.Role}"),
    };

        var jwtKey = _configuration["Jwt:Key"] ?? _configuration["JWT:KEY"];
        if (string.IsNullOrWhiteSpace(jwtKey))
        {
            throw new InvalidOperationException("JWT key is not configured.");
        }

        var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtKey));
        var credentials = new SigningCredentials(key, SecurityAlgorithms.HmacSha256);
        var issuer = _configuration["Jwt:Issuer"] ?? "EBookStore";
        var audience = _configuration["Jwt:Audience"] ?? "EBookStoreUsers";
        var expiryMinutes = int.TryParse(_configuration["Jwt:ExpiryMinutes"], out var parsedMinutes)
            ? parsedMinutes
            : 60;

        var token = new JwtSecurityToken(
              issuer: issuer,
              audience: audience,
              claims: claims,
              expires: DateTime.UtcNow.AddMinutes(expiryMinutes),
              signingCredentials: credentials
          );

        return new JwtSecurityTokenHandler().WriteToken(token);
    }
}