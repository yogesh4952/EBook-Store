namespace EbookStore.Controllers.Authenticaiton;

using EbookStore.Dtos.RegisterDto;
using EbookStore.Entities.User;
using EbookStore.Services;
using Microsoft.AspNetCore.Identity;
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
    public ActionResult<User> register(RegisterDto request)
    {
        var hashedPassword = new PasswordHasher<User>().HashPassword(user, request.Password);
        request.Password = hashedPassword;

        request.Password = hashedPassword;
        return Ok(request);
    }

}