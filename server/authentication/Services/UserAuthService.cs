using EbookStore.Data;
using EbookStore.Dtos.LoginDto;
using EbookStore.Dtos.RegisterDto;
using EbookStore.Entities.User;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;

namespace EbookStore.Services;

public class AuthService
{
    private readonly AppDbContext _context;

    public AuthService(AppDbContext context)
    {
        _context = context;
    }

    public async Task<User> Register(RegisterDto dto)
    {
        var exisitingUsers = await _context.Users.FirstOrDefaultAsync(x => x.Email == dto.Email);

        if (exisitingUsers != null)
        {
            throw new Exception("Email already exists!");
        }

        var hasher = new PasswordHasher<User>();
        var user = new User();
        user.PasswordHash = hasher.HashPassword(user, dto.Password);
        user.Email = dto.Email;
        user.FirstName = dto.FirstName;
        user.LastName = dto.LastName;

        _context.Users.Add(user);

        await _context.SaveChangesAsync();
        return user;
    }

    public async Task<User> Login(LoginDto dto)
    {
        var hasher = new PasswordHasher<User>();

        var user = await _context.Users
            .FirstOrDefaultAsync(u => u.Email == dto.Email);

        if (user == null)
            throw new BadHttpRequestException("Email doesn't exist.");

        var result = hasher.VerifyHashedPassword(
            user,
            user.PasswordHash,
            dto.Password
        );

        if (result == PasswordVerificationResult.Failed)
            throw new BadHttpRequestException("Incorrect password.");

        return user;
    }
}
