package types

var UserInviteEMailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Vapusdata {Account} Login</title>
</head>
<body>
	<div style="background-color: #f5f5f5; padding: 20px;">
		<div style="background-color: #fff; padding: 20px; border-radius: 10px;">
			<h1 style="color: #333; font-size: 24px; font-weight: 600; margin-bottom: 20px;">Welcome to VapusData Platform</h1>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">Hi {Name},</p>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">You have been invited to join VapusData. Please click the link below to login and access the platform.</p>
			<a href="{Link}" style="background-color: #007bff; color: #fff; padding: 10px 20px; border-radius: 5px; text-decoration: none; display: inline-block; margin-bottom: 20px;">Login</a>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">Thanks,</p>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">{Account}</p>
		</div>
	</div>
</body>
</html>
`

var DataProductAccessRequestTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Vapusdata {Account} Data Product Access Request</title>
</head>
<body>
	<div style="background-color: #f5f5f5; padding: 20px;">
		<div style="background-color: #fff; padding: 20px; border-radius: 10px;">
			<h1 style="color: #333; font-size: 24px; font-weight: 600; margin-bottom: 20px;">Data Product Access Request</h1>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">Hi,</p>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">User {UserId} has requested data product access from ORGANIZATION {ORGANIZATION} for data product {DataProduct}. Please click the link below to view the request.</p>
			<a href="{Link}" style="background-color: #007bff; color: #fff; padding: 10px 20px; border-radius: 5px; text-decoration: none; display: inline-block; margin-bottom: 20px;">Data Product Access Request</a>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">Thanks,</p>
			<p style="color: #333; font-size: 16px; margin-bottom: 20px;">{Account}</p>
		</div>
	</div>
</body>
</html>
`
