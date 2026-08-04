namespace EbookStore.Controllers.Authenticaiton;

using EbookStore.Dtos.RegisterDto;
using EbookStore.Entities.User;
using EbookStore.Services;
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

}