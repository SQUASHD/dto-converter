using System;
using System.ComponentModel.DataAnnotations;
using StudentBlogAPI.Validation;

namespace StudentBlogAPI.Models.DTOs;

public record UserProfileResDto(
    Guid UserId,
    string DisplayName,
    DateTime DateOfBirth,
    [EmailAddress] string Email,
    string Bio
);

public record ArticleCreateReqDto(
    [Required] string Title,
    [Required] [MaxLength(5000)] string Content,
    [Url] string FeaturedImageUrl,
    Guid AuthorId
);

public record CommentPostReqDto(
    Guid ArticleId,
    Guid UserId,
    [Required] [MaxLength(1000)] string Body
);

public record EnrollmentReqDto(
    Guid StudentId,
    Guid CourseId,
    [Range(1, 4)] [Required] int Year
);

public record CourseDetailResDto(
    Guid CourseId,
    string CourseName,
    string Description,
    [Range(1, 8)] int Semester,
    string ProfessorName
);

public record PaymentProcessReqDto(
    Guid UserId,
    [CreditCard] string CardNumber,
    DateTime CardExpiry,
    [Range(100, 999)] int CVV
);

public record EventScheduleReqDto(
    [Required] string EventTitle,
    DateTime StartTime,
    DateTime EndTime,
    [MaxLength(1000)] string Description,
    [EmailAddress] string OrganizerEmail
);

public record GradeUpdateReqDto(
    Guid EnrollmentId,
    [Range(0.0, 100.0)] double Score
);

public record BookBorrowReqDto(
    Guid UserId,
    Guid BookId,
    [Required] DateTime BorrowDate,
    DateTime ExpectedReturnDate
);
