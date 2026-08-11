namespace EbookStore.Controllers.Authenticaiton;

using System.Security.Claims;
using EbookStore.Dtos.LoginDto;
using EbookStore.Dtos.RegisterDto;
using EbookStore.Entities.User;
using EbookStore.Services;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;


[Route("api/[controller]")]
[ApiController]
public class AuthController : ControllerBase
{
    private readonly AuthService _authService;

    public AuthController(AuthService authService)
    {
        _authService = authService;
    }


    [HttpPost("register")]
    public async Task<ActionResult<User>> Register(RegisterDto dto)
    {
        var user = await _authService.Register(dto);

        return Ok(user);
    }


    [HttpPost("login")]
    public async Task<ActionResult<LoginResponseDto>> Login(LoginDto dto)
    {
        var user = await _authService.Login(dto);
        return Ok(user);
    }

    [HttpGet("me")]
    [Authorize]
    public async Task<IActionResult> Me()
    {
        var userEmail = User.FindFirst(ClaimTypes.Email)?.Value
                     ?? User.FindFirst("email")?.Value;

        if (string.IsNullOrEmpty(userEmail))
        {
            return Unauthorized(new { message = "Invalid token claims." });
        }

        var user = await _authService.GetByEmail(userEmail);

        if (user == null)
        {
            return NotFound(new { message = "User not found." });
        }

        // Return safe user details (excluding PasswordHash)
        return Ok(new
        {
            user.Id,
            user.Email,
            user.FirstName,
            user.LastName
        });
    }

}