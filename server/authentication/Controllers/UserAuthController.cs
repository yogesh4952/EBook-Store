namespace EbookStore.Controllers.Authenticaiton;

using EbookStore.Dtos.LoginDto;
using EbookStore.Entities.User;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;


[Route("api/[controller]")]
[ApiController]
public class AuthController : ControllerBase
{

    public static User user = new();
    [HttpPost("register")]
    public ActionResult<User> register(UserDto request)
    {
        var hashedPassword = new PasswordHasher<User>().HashPassword(user, request.Password);
        request.Password = hashedPassword;

        request.Password = hashedPassword;
        return Ok(request);
    }

}