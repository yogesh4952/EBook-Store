using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace EbookStore.Dtos.RegisterDto;

public class RegisterDto
{


    [Required(ErrorMessage = "FirstName is required.")]

    public required string FirstName { get; set; }

    [Required(ErrorMessage = "LastName is required.")]

    public required string LastName { get; set; }


    [Required(ErrorMessage = "Email is required.")]


    public string Email { get; set; }


    [Required(ErrorMessage = "Password is required.")]

    public string Password { get; set; }


    public int Role { get; set; } = 0;
}