using System.ComponentModel.DataAnnotations;

namespace EbookStore.Dtos.RegisterDto;

public class RegisterDto
{


    [Required(ErrorMessage = "FirstName is required.")]

    public required string FirstName { get; set; }

    [Required(ErrorMessage = "LastName is required.")]
    public required string LastName { get; set; }


    [Required(ErrorMessage = "Email is required.")]
    [EmailAddress(ErrorMessage = "Invalid email address.")]
    public string Email { get; set; } = string.Empty;


    [Required(ErrorMessage = "Password is required.")]
    public string Password { get; set; } = string.Empty;


    public int Role { get; set; } = 0;
}