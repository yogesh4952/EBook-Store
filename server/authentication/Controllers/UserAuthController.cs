namespace EbookStore.Controllers.Authenticaiton;

using EbookStore.Entities.U;
using Microsoft.AspNetCore.Components;
using Microsoft.AspNetCore.Mvc;


[Route("api/[controller]")]
[ApiController]
public class AuthController : ControllerBase
{

    [HttpPost("register")]
    public ActionResult<User> register()
    {

    }

}