package verification

import (
	"io"

	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	conv "REDACTED/team-11/backend/admin/internal/controller/http/v1/converter"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/validator"
	resp "REDACTED/team-11/backend/admin/internal/controller/response"
	"REDACTED/team-11/backend/admin/internal/entity"
	ct "REDACTED/team-11/backend/admin/pkg/utils/controller"
)

type Verification struct {
	usecase VerificationUseCase
}

func New(uc VerificationUseCase) *Verification {
	return &Verification{
		usecase: uc,
	}
}

// @Summary Check verification
// @Description Returns user verification data. Only for ADMINs
// @Tags Verification
// @Produce json
// @Security Bearer
// @Param        id    path     string  true  "User id"  Format(uuid)
// @Success 200 {object} dto.VerificationData "Successful response"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/verification/{id}/check [get]
func (v *Verification) CheckVerify(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	data, err := v.usecase.CheckVerify(ctx, id)
	if err != nil {
		v, ok := err.GetTag("coffee").(bool)
		if ok && v {
			resp.AbortErrMsg(c, err)
			return
		} else if err.GetCode() == e.NotFound {
			c.JSON(httper.StatusOK, conv.DtoVerify(false, ""))
			return
		} else {
			resp.AbortErrMsg(c, err)
			return
		}
	}

	result := conv.DtoVerify(true, data.PassportImage)

	c.JSON(httper.StatusOK, result)
}

// @Summary Set verification
// @Description Set user verification data. Only for ADMINs
// @Tags Verification
// @Security Bearer
// @Accept multipart/form-data
// @Param id path string true "User id" Format(uuid)
// @Param passport formData file false "passport image"
// @Success 204 "Successful setup of verification data."
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/verification/{id}/set [post]
func (v *Verification) Verify(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	file, err := c.FormFile("passport")
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}

	reader, err := file.Open()
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}
	defer reader.Close()

	buffer, err := io.ReadAll(reader)
	if err != nil {
		resp.AbortErrMsg(c, badReqErr)
		return
	}

	image := &entity.Image{
		Name:        file.Filename,
		Buffer:      buffer,
		Size:        file.Size,
		ContentType: file.Header["Content-Type"][0],
	}

	if err := v.usecase.Verify(ctx, id, image); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}
