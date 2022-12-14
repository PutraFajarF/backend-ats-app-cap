package main

import (
	"github.com/PutraFajarF/backend-ats-app-cap/config"
	"github.com/PutraFajarF/backend-ats-app-cap/controller"
	"github.com/PutraFajarF/backend-ats-app-cap/middleware"
	"github.com/PutraFajarF/backend-ats-app-cap/repository"
	"github.com/PutraFajarF/backend-ats-app-cap/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB = config.SetupConnectionDatabase()

	userRepository           repository.UserRepository           = repository.NewUserRepository(db)
	applicantRepository      repository.ApplicantRepository      = repository.NewApplicantRepository(db)
	employeeRepository       repository.EmployeeRepository       = repository.NewEmployeeRepository(db)
	experienceRepository     repository.ExperienceRepository     = repository.NewExperienceRepository(db)
	skillrepository          repository.SkillRepository          = repository.NewSkillRepository(db)
	jobsRepository           repository.JobsRepository           = repository.NewJobsRepository(db)
	jobApplicationRepository repository.JobApplicationRepository = repository.NewjobApplication(db)
	jobAppliedRepository     repository.JobsAppliedRepository    = repository.NewJobAppliedConnection(db)
	jobApplicationEmployeeRepository     repository.JobApplicationEmployeeRepository    = repository.NewjobApplicationEmployee(db)

	jwtService service.JWTService = service.NewJWTService()

	authService       service.AuthService           = service.NewAuthService(userRepository, applicantRepository, employeeRepository)
	applicantService  service.ApplicantService      = service.NewApplicantService(applicantRepository)
	employeeService   service.EmployeeService       = service.NewEmployeeService(employeeRepository)
	experienceService service.ExperienceService     = service.NewExperienceService(experienceRepository)
	skillService      service.SkillService          = service.NewSkillService(skillrepository)
	jobsService       service.JobsService           = service.NewJobsService(jobsRepository)
	jobApplication    service.JobApplicationService = service.NewJobApplicationService(jobApplicationRepository)
	jobAppliedService service.JobAppliedService     = service.NewJobAppliedService(jobAppliedRepository)
	jobApplicationServiceEmployee service.JobApplicationEmployeeService     = service.NewJobApplicationEmployeeService(jobApplicationEmployeeRepository)

	authController       controller.AuthController       = controller.NewAuthController(authService, jwtService)
	applicantController  controller.ApplicantController  = controller.NewApplicantController(applicantService, jwtService)
	employeeController   controller.EmployeeController   = controller.NewEmployeeController(employeeService, jwtService)
	experienceController controller.ExperienceController = controller.NewExperienceController(experienceService, jwtService)
	skillController      controller.SkillController      = controller.NewSkillController(skillService, jwtService)
	jobsController       controller.JobsController       = controller.NewJobsController(jobsService, jwtService)
	jobsAppController    controller.JobsAppController    = controller.NewJobsAppController(jobsService)
	jobAppliedController controller.JobAppliedController = controller.NewJobAppliedController(jobAppliedService, jwtService)
	jobApplicantion controller.JobApplicationController = controller.NewJobApplicationController(jobApplication, jwtService)
	jobApplicantionEmployee controller.JobApplicationEmployeeController = controller.NewJobApplicationEmployeeController(jobApplicationServiceEmployee, jwtService)
)

func main() {
	defer config.CloseConnectionDatabase(db)
	r := gin.Default()
	r.Static("/images", config.ProjectRootPath+"/images/applicants/")
	// r.Static("/images", "./images")
	r.Use(cors.Default())
	authRoutes := r.Group("/")
	{
		authRoutes.POST("login", authController.Login)
		authRoutes.POST("applicants/register", authController.RegisterApplicants)
		authRoutes.GET("/jobs", jobsAppController.GetAllJobsApplicant)
		authRoutes.GET("/jobs/:id", jobsAppController.ApplicantGetJobsByID)
	}

	// Applicant Routes
	authApplicantRoutes := r.Group("/applicants", middleware.AuthorizeJWT(jwtService))
	{
		authApplicantRoutes.PUT("/users/:id", applicantController.EditApplicant)
		authApplicantRoutes.GET("/users/fetch", applicantController.FetchUserApplicant)
		authApplicantRoutes.POST("/avatars", applicantController.UploadAvatar)

		authApplicantRoutes.GET("/cv", applicantController.DetailApplicant)

		authApplicantRoutes.POST("/jobexperiences", experienceController.CreateExperience)
		authApplicantRoutes.PUT("/jobexperiences/:id", experienceController.UpdateExperience)
		authApplicantRoutes.DELETE("/jobexperiences/:id", experienceController.DeleteExperience)
		authApplicantRoutes.GET("/jobexperiences/", experienceController.GetAllExperiences)
		authApplicantRoutes.GET("/jobexperiences/:id", experienceController.GetExperienceByID)

		authApplicantRoutes.POST("/skills", skillController.CreateSkill)
		authApplicantRoutes.PUT("/skills/:id", skillController.UpdateSkill)
		authApplicantRoutes.GET("/skills/:id", skillController.GetSkillByID)
		authApplicantRoutes.DELETE("/skills/:id", skillController.DeleteSkill)
		authApplicantRoutes.GET("/skills", skillController.GetSkills)

		authApplicantRoutes.POST("/jobapplication", jobApplicantion.CreateApply)

		authApplicantRoutes.GET("/applied", jobAppliedController.JobsAppliedByApplicantID)

	}

	// Employees Routes
	authEmployeeRoutes := r.Group("/employees", middleware.AuthorizeJWT(jwtService))
	{
		authEmployeeRoutes.POST("/register", authController.RegisterEmployees)
		authEmployeeRoutes.PUT("/users/:id", employeeController.EditEmployee)
		authEmployeeRoutes.GET("/users/fetch", employeeController.FetchUserEmployee)

		authEmployeeRoutes.GET("/jobs", jobsController.GetAllJobs)
		authEmployeeRoutes.GET("/jobs/:id", jobsController.GetJobsByID)
		authEmployeeRoutes.POST("/jobs", jobsController.CreatedJobs)
		authEmployeeRoutes.DELETE("/jobs/:id", jobsController.DeleteJobs)
		authEmployeeRoutes.PUT("/jobs/:id", jobsController.UpdateJobs)

		authEmployeeRoutes.PUT("/progress/:id", jobApplicantionEmployee.ProgressApplication)
	}
	r.Run()
}
