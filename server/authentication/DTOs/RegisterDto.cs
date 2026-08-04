using System.ComponentModel.DataAnnotations;

namespace EbookStore.Dtos.RegisterDto;

public class RegisterDto
{


    [Required(ErrorMessage = "FirstName is required.")]

    public string FirstName { get; set; } = string.Empty;

    [Required(ErrorMessage = "LastName is required.")]

    public string LastName { get; set; } = string.Empty;


    [Required(ErrorMessage = "Email is required.")]


    public string Email { get; set; } = string.Empty;


    [Required(ErrorMessage = "Password is required.")]

    public string Password { get; set; } = string.Empty;
}